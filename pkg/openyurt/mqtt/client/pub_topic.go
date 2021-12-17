package client

type PublishTopicor interface {
	GetPublishPreTopic() string
	GetPublishGetTopic(name string) string
	GetPublishCreateTopic(name string) string
	GetPublishUpdateTopic(name string) string
	GetPublishPatchTopic(name string) string
	GetPublishDeleteTopic(name string) string
}
