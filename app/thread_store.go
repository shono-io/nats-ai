package app

import "github.com/henomis/lingoose/thread"

type ThreadStore interface {
	GetThread(threadID string) (*thread.Thread, error)
	StoreThread(threadID string, thread *thread.Thread) error
}
