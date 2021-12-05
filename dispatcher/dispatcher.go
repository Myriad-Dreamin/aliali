package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/database"
	"github.com/Myriad-Dreamin/aliali/dispatcher/service"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gorm.io/gorm"
	"log"
	"os"
)

type ServiceContext struct {
	S         suppress.ISuppress
	authedAli *ali_drive.Ali
	Impl      service.IService
}

type Dispatcher struct {
	configPath string
	dbPath     string

	dbMgr       *DBManager
	cfgMgr      *ConfigManager
	cfg         *ali_notifier.Config
	auth        *model.AliAuthModel
	httpHeaders [][2]string

	s         suppress.ISuppress
	ali       *ali_drive.Ali
	authedAli *ali_drive.Ali
	db        *gorm.DB
	xdb       *database.DB
	notifier  ali_notifier.INotifier
	logger    *log.Logger

	fileUploads  chan *ali_notifier.FsUploadRequest
	serviceQueue chan *ServiceContext
	tokenInvalid chan bool
}

type Option = func(w *Dispatcher) *Dispatcher

func MockDB() Option {
	return func(w *Dispatcher) *Dispatcher {
		w.db = w.openMock()
		if w.db == nil {
			return nil
		}
		return w
	}
}

func WithConfigPath(cfgPath string) Option {
	return func(w *Dispatcher) *Dispatcher {
		w.configPath = cfgPath
		return w
	}
}

func WithDBPath(dbPath string) Option {
	return func(w *Dispatcher) *Dispatcher {
		w.dbPath = dbPath
		return w
	}
}

func WithConfig(cfg *ali_notifier.Config) Option {
	return func(w *Dispatcher) *Dispatcher {
		w.cfg = cfg
		return w
	}
}

func WithNotifier(notifier ali_notifier.INotifier) Option {
	return func(w *Dispatcher) *Dispatcher {
		w.notifier = notifier
		return w
	}
}

func WithServiceReplicate(services ...service.IService) Option {
	return func(w *Dispatcher) *Dispatcher {
		w.serviceQueue = make(chan *ServiceContext, len(services))
		for i := range services {
			w.serviceQueue <- &ServiceContext{Impl: services[i]}
		}
		return w
	}
}

func NewDispatcher(options ...Option) *Dispatcher {
	s := suppress.PanicAll{}

	var httpHeaders [][2]string
	httpHeaders = append(httpHeaders, [2]string{"origin", "https://aliyundrive.com"})
	httpHeaders = append(httpHeaders, [2]string{"referer", "https://aliyundrive.com"})

	var w = &Dispatcher{
		s:           s,
		fileUploads: make(chan *ali_notifier.FsUploadRequest, 10),
		httpHeaders: httpHeaders,
	}

	w.dbMgr = &DBManager{S: w.s}
	w.cfgMgr = &ConfigManager{S: w.s}
	for _, option := range options {
		w = option(w)
		if w == nil {
			return nil
		}
	}

	w.ali = w.makeAliClient(w.s)
	if w.cfg == nil {
		if len(w.configPath) == 0 {
			w.configPath = DefaultConfigPath
		}
		w.cfg = w.syncConfig()
	}
	if w.logger == nil {
		w.logger = log.New(os.Stderr, "[dispatcher] ", log.LUTC|log.Llongfile)
	}
	if w.db == nil {
		if len(w.dbPath) == 0 {
			w.dbPath = DefaultDatabasePath
		}
		w.db = w.openDB()
	}
	w.tokenInvalid = make(chan bool, 1)
	if w.serviceQueue == nil {
		w.serviceQueue = make(chan *ServiceContext, 1)
		w.serviceQueue <- &ServiceContext{
			Impl: &service.UploadImpl{
				Logger: log.New(os.Stderr, "[service] ", log.LUTC|log.Llongfile),
			},
		}
	}
	if w.xdb == nil {
		w.xdb = &database.DB{ISuppress: w.s}
	}
	if w.notifier == nil {
		w.notifier = &ali_notifier.HttpRecorderNotifier{
			CapturePath: "/rec",
		}
	}

	if w.ali == nil || w.db == nil || w.cfg == nil || w.logger == nil ||
		w.serviceQueue == nil || w.xdb == nil || w.notifier == nil {
		return nil
	}

	w.setupNotifier()

	return w
}
