package adapters

import (
	"github.com/rbaderts/poker/domain"
	"sync"
)


func init() {
}

type NullDispatcher struct {
	seatChannels     map[string]map[int]string
	blockedRequestChannels map[int]*chan domain.Message

	lock sync.Mutex
}


func NewNullDispatcher() domain.MessageDispatcher {

	d := new(NullDispatcher)
	return d
}

func (this *NullDispatcher) RespondToRequest(requestId int, msg *domain.Message) error {

	return nil
}
func (this *NullDispatcher) PublishToTable(table *domain.Table, data interface{}) error {return nil}

func (this *NullDispatcher) BroadcastToTable(
	table *domain.Table, typ domain.EventType, data interface{})  error {
	return nil
}


func (this *NullDispatcher) SendToSeat(
	table *domain.Table, seatnum int, typ domain.EventType, data interface{}) (*domain.Event, error) {
		return nil, nil
}

func (this *NullDispatcher) SendRequestAndWait(
	table *domain.Table,
	seatnum int,
	typ domain.EventType,
	data interface{},
	handler domain.ResponseHandler) (interface{}, error) {
		return nil, nil
}
