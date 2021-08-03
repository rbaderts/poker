package adapters

import (
	"context"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/nats-io/nats.go"
	"github.com/rbaderts/poker/domain"
	"sync"
	_ "time"
)

var NatsClient *nats.Conn

func init() {

	//	var err error
	//	NatsClient, err = nats.Connect(nats.DefaultURL)

	//	if err != nil {
	//		fmt.Printf("Error: %v\n", err)
	//	}

}

type NatsDispatcher struct {
	seatChannels           map[string]map[int]string
	blockedRequestChannels map[int]*chan domain.Message

	lock sync.Mutex
}

func NewNatsDispatcher() domain.MessageDispatcher {

	d := new(NatsDispatcher)
	d.seatChannels = make(map[string]map[int]string, 0)
	d.blockedRequestChannels = make(map[int]*chan domain.Message)

	return d
}

func (this *NatsDispatcher) RespondToRequest(requestId int, msg *domain.Message) error {

	fmt.Printf("Received response to request %d\n", requestId)
	waitChan := this.blockedRequestChannels[requestId]

	if waitChan != nil {
		*waitChan <- *msg

		this.lock.Lock()
		delete(this.blockedRequestChannels, requestId)
		this.lock.Unlock()
	} else {
		fmt.Printf("No channel found waiting for resp to requestId %d\n", requestId)
	}

	return nil
}

func (this *NatsDispatcher) PublishToTable(
	table *domain.Table, data interface{}) error {
	ch := table.Channel
	err := NatsClient.Publish(ch, data.([]byte))

	if err != nil {
		fmt.Printf("Error publishing msg: %v\n", err)
	}
	return nil
}

func (this *NatsDispatcher) BroadcastToTable(
	table *domain.Table, typ domain.EventType, data interface{}) error {

	ch := table.Channel
	/*
		ch, ok := this.tableChannels[table.Id]
		if !ok {
			ch = "t_" + table.Id[0:4]
			this.tableChannels[table.Id] = ch
		}
	*/

	event := domain.NewEvent(typ, data)

	eventData, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling json: %v\n", err)
		return err
	}

	blue := color.New(color.FgBlue).SprintFunc()
	color.Set(color.FgRed)
	fmt.Printf("Publishing event type %s to channel %s\n", blue(typ.String()), ch)
	color.Unset()
	//	_, err = CentClient.Publish(this.Channel, data)
	//ctx := context.Background()

	err = NatsClient.Publish(ch, eventData)

	if err != nil {
		fmt.Printf("Error publishing msg: %v\n", err)
		return err
	}

	return nil
}

func (this *NatsDispatcher) SendToSeat(
	table *domain.Table, seatnum int, typ domain.EventType, data interface{}) (*domain.Event, error) {

	ch := table.Seats[seatnum].Channel

	event := domain.NewEvent(typ, data)

	eventData, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling json: %v\n", err)
		return nil, err
	}

	blue := color.New(color.FgBlue).SprintFunc()
	color.Set(color.FgRed)
	fmt.Printf("Publishing event type %s to channel %s\n", blue(typ.String()), ch)
	color.Unset()
	ctx := context.Background()
	_ = ctx

	err = NatsClient.Publish(ch, eventData)

	if err != nil {
		fmt.Printf("Error publishing msg: %v\n", err)
		return nil, err
	}
	return event, nil
}

func (this *NatsDispatcher) SendRequestAndWait(
	table *domain.Table,
	seatnum int,
	typ domain.EventType,
	data interface{},
	handler domain.ResponseHandler) (interface{}, error) {

	ch := table.Seats[seatnum].Channel
	_ = ch

	//	NatsClient.Request(ch, data, )
	//	func (nc *Conn) Request(subj string, data []byte, timeout time.Duration) (*Msg, error) {

	//		event, err := this.SendToSeat(table, seatnum, typ, data)
	//	_ = err

	/*
		waitChan := make(chan domain.Message)

		id := event.Id
		this.lock.Lock()
		this.blockedRequestChannels[id] = &waitChan
		this.lock.Unlock()

		defer func() {
			this.lock.Lock()
			this.blockedRequestChannels[id] = nil
			this.lock.Unlock()
		}()

		select {
		case res := <-waitChan:

			//		amount := res.Data.(BetResponse).Amount
			//		amountStr := res.Data.(BetData)
			///		amountStr := res.Params[AMOUNT]
			//		amount, _ := strconv.Atoi(amountStr)
			result, err := handler.Handle(&res)
			/*
				err := this.PlayerBet(int(amount), seat)
	*/

	/*
			if err != nil {
				fmt.Errorf("error recording player bet: %v\n", err)
				return 0, err
			}

			return result, nil

		case <-time.After(300 * time.Second):
			fmt.Println("timeout after 5 minutes")
			return 0, errors.New("Bet Reqeust timeout")
		}


	*/
	return nil, nil

}
