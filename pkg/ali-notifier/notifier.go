package ali_notifier

type ISubscription interface {
	OnFsUpload(req *FsUploadRequest)
}

type INotifier interface {
	RegisterCallback(s ISubscription)
}

type Notifier struct {
	subscriptions []ISubscription
}

func (n *Notifier) RegisterCallback(s ISubscription) {
	n.subscriptions = append(n.subscriptions, s)
}

func (n *Notifier) Emit(req *FsUploadRequest) {
	for i := range n.subscriptions {
		n.subscriptions[i].OnFsUpload(req)
	}
}
