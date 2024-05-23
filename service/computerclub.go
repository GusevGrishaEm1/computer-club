// Пакет service предоставляет функциональность для управления операциями компьютерного клуба,
// включая обработку событий клиентов, управление столами и расчет стоимости использования.
// Этот пакет обрабатывает события и поддерживает состояние клуба, такое как распределение
// клиентов по столам, количество свободных столов и журналы событий.
package service

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// INPUT_ACTION определяет типы входных действий, которые может выполнять клиент.
type INPUT_ACTION int

// OUTPUT_ACTION определяет типы выходных действий, возникающих в результате действий клиента.
type OUTPUT_ACTION int

// ERROR_MESSAGE определяет сообщения об ошибках для различных недопустимых операций.
type ERROR_MESSAGE string

// client представляет клиента в компьютерном клубе.
type client string

// table представляет стол в компьютерном клубе.
type table int

const (
	// CLIENT_IN указывает на вход клиента в клуб.
	CLIENT_IN INPUT_ACTION = iota + 1
	// CLIENT_TABLE указывает на запрос клиента на стол.
	CLIENT_TABLE
	// CLIENT_WAIT указывает на ожидание клиента в очереди.
	CLIENT_WAIT
	// CLIENT_OUT указывает на выход клиента из клуба.
	CLIENT_OUT
)

const (
	// CLIENT_EXIT указывает на то, что клиент покидает клуб.
	CLIENT_EXIT OUTPUT_ACTION = iota + 11
	// CLIENT_POP_QUEUE указывает на то, что клиент удален из очереди.
	CLIENT_POP_QUEUE
	// ERROR указывает на то, что произошла ошибка.
	ERROR
)

const (
	// YOU_SHALL_NOT_PASS указывает на то, что клиент уже находится в клубе.
	YOU_SHALL_NOT_PASS = ERROR_MESSAGE("YouShallNotPass")
	// CLIENT_UNKNOWN указывает на неизвестного клиента.
	CLIENT_UNKNOWN = ERROR_MESSAGE("ClientUnknown")
	// PLACE_IS_BUSY указывает на занятое место.
	PLACE_IS_BUSY = ERROR_MESSAGE("PlaceIsBusy")
	// NOT_OPEN_YET указывает на то, что клуб еще не открыт.
	NOT_OPEN_YET = ERROR_MESSAGE("NotOpenYet")
	// I_CAN_WAIT_NO_LONGER указывает на то, что клиент больше не может ждать.
	I_CAN_WAIT_NO_LONGER = ERROR_MESSAGE("ICanWaitNoLonger")
)

// event представляет событие, которое происходит в клубе.
type event struct {
	time   time.Time
	action INPUT_ACTION
	client client
	table  table
	source string
}

// computerClub представляет состояние и логику работы компьютерного клуба.
type computerClub struct {
	clients            map[client]table
	tablesAndStartTime map[table]time.Time
	tablesAndHours     map[table]time.Duration
	events             []event
	startTime          time.Time
	endTime            time.Time
	costOfHour         int
	freeTables         int
	countOfTables      int
	queueOfClients     *Queue
	outputs            []string
	outputsAfterEnd    []string
}

// NewComputerClub создает новый экземпляр компьютерного клуба.
func NewComputerClub() *computerClub {
	return &computerClub{}
}

// SetInput загружает конфигурацию и события из файла.
func (c *computerClub) SetInput(config Config) error {
	file, err := os.Open(config.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := c.parseConfig(scanner); err != nil {
		return err
	}

	events, err := c.parseEvents(scanner)
	if err != nil {
		return err
	}
	c.events = events
	return nil
}

// parseConfig парсит конфигурационные данные из файла.
func (c *computerClub) parseConfig(scanner *bufio.Scanner) error {
	if !scanner.Scan() {
		return fmt.Errorf("failed to read number of tables")
	}
	numTables, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("failed to parse number of tables: %v %s", err, scanner.Text())
	}

	if !scanner.Scan() {
		return fmt.Errorf("failed to read start and end time of the day")
	}
	parts := strings.Split(scanner.Text(), " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid start and end time format %s", scanner.Text())
	}
	start, err := time.Parse("15:04", parts[0])
	if err != nil {
		return fmt.Errorf("failed to parse start time: %v", err)
	}
	end, err := time.Parse("15:04", parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse end time: %v", err)
	}

	if !scanner.Scan() {
		return fmt.Errorf("failed to read cost of an hour")
	}
	costOfHour, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("failed to parse cost of an hour: %v %s", err, scanner.Text())
	}

	c.tablesAndStartTime = make(map[table]time.Time, numTables)
	c.countOfTables = numTables
	c.freeTables = numTables
	c.startTime = start
	c.endTime = end
	c.costOfHour = costOfHour
	c.clients = make(map[client]table)
	c.tablesAndHours = make(map[table]time.Duration, numTables)
	c.queueOfClients = NewQueue()

	return nil
}

// parseEvents парсит события из файла.
func (c *computerClub) parseEvents(scanner *bufio.Scanner) ([]event, error) {
	var events []event
	for scanner.Scan() {
		lineEvent := scanner.Text()
		parts := strings.Split(lineEvent, " ")
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid event format %s", lineEvent)
		}
		time, err := time.Parse("15:04", parts[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %v", err)
		}
		action, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse action: %v", err)
		}
		clientName := parts[2]
		var tableNum int
		if len(parts) == 4 {
			tableNum, err = strconv.Atoi(parts[3])
			if err != nil {
				return nil, fmt.Errorf("failed to parse table: %v", err)
			}
		}
		events = append(events, event{
			time:   time,
			action: INPUT_ACTION(action),
			client: client(clientName),
			table:  table(tableNum),
			source: lineEvent,
		})
	}
	return events, nil
}

// ProcessInput обрабатывает входные данные и события, обновляя состояние клуба.
func (c *computerClub) ProcessInput() error {
	c.initOutputs()
	for _, event := range c.events {
		if !c.filterByTime(event) {
			continue
		}
		c.outputs = append(c.outputs, event.source)
		c.processEvent(event)
	}
	c.processEndOfDay()
	return nil
}

// initOutputs инициализирует выходные данные.
func (c *computerClub) initOutputs() {
	c.outputs = []string{fmt.Sprintf("%02d:%02d", c.startTime.Hour(), c.startTime.Minute())}
	c.outputsAfterEnd = []string{}
}

// filterByTime фильтрует события по времени.
func (c *computerClub) filterByTime(event event) bool {
	if c.startTime.After(event.time) {
		c.outputs = append(c.outputs, event.source, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, NOT_OPEN_YET))
		return false
	}
	if c.endTime.Before(event.time) {
		c.outputsAfterEnd = append(c.outputsAfterEnd, event.source, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, NOT_OPEN_YET))
		return false
	}
	return true
}

// processEvent обрабатывает одно событие.
func (c *computerClub) processEvent(event event) {
	switch event.action {
	case CLIENT_IN:
		c.handleClientIn(event)
	case CLIENT_TABLE:
		c.handleClientTable(event)
	case CLIENT_WAIT:
		c.handleClientWait(event)
	case CLIENT_OUT:
		c.handleClientOut(event)
	}
}

// handleClientIn обрабатывает вход клиента в клуб.
func (c *computerClub) handleClientIn(event event) {
	if _, ok := c.clients[event.client]; ok {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, YOU_SHALL_NOT_PASS))
		return
	}
	c.clients[event.client] = 0
}

// handleClientTable обрабатывает запрос клиента на стол.
func (c *computerClub) handleClientTable(event event) {
	if _, ok := c.clients[event.client]; !ok {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, CLIENT_UNKNOWN))
		return
	}
	if _, ok := c.tablesAndStartTime[event.table]; ok {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, PLACE_IS_BUSY))
		return
	}
	c.tablesAndStartTime[event.table] = event.time
	c.clients[event.client] = event.table
	c.freeTables--
}

// handleClientWait обрабатывает ожидание клиента в очереди.
func (c *computerClub) handleClientWait(event event) {
	if c.freeTables > 0 {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, I_CAN_WAIT_NO_LONGER))
		return
	}
	if c.queueOfClients.Size == c.countOfTables {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), CLIENT_EXIT, event.client))
		return
	}
	c.queueOfClients.Enqueue(event.client)
}

// handleClientOut обрабатывает выход клиента из клуба.
func (c *computerClub) handleClientOut(event event) {
	val, ok := c.clients[event.client]
	if !ok {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", event.time.Hour(), event.time.Minute(), ERROR, CLIENT_UNKNOWN))
		return
	}
	c.updateTableUsage(val, event.time)
	delete(c.tablesAndStartTime, val)
	c.freeTables++
	c.processQueue(event, val)
	delete(c.clients, event.client)
}

// updateTableUsage обновляет данные о времени использования стола.
func (c *computerClub) updateTableUsage(val table, currentTime time.Time) {
	startTime, ok := c.tablesAndStartTime[val]
	if ok {
		diff := currentTime.Sub(startTime)
		c.tablesAndHours[val] += diff
	}
}

// processQueue обрабатывает очередь клиентов.
func (c *computerClub) processQueue(event event, val table) {
	if c.queueOfClients.Size > 0 {
		client, ok := c.queueOfClients.Dequeue()
		if ok {
			c.clients[client] = val
			c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s %d", event.time.Hour(), event.time.Minute(), CLIENT_POP_QUEUE, client, val))
			c.tablesAndStartTime[val] = event.time
			c.clients[client] = val
			c.freeTables--
		}
	}
}

// processEndOfDay обрабатывает конец рабочего дня клуба.
func (c *computerClub) processEndOfDay() {
	sortedClients := make([]client, 0, len(c.clients))
	for client := range c.clients {
		sortedClients = append(sortedClients, client)
	}
	sort.Slice(sortedClients, func(i, j int) bool {
		return string(sortedClients[i]) < string(sortedClients[j])
	})
	for _, client := range sortedClients {
		c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d %d %s", c.endTime.Hour(), c.endTime.Minute(), CLIENT_EXIT, client))
	}
	c.outputs = append(c.outputs, fmt.Sprintf("%02d:%02d", c.endTime.Hour(), c.endTime.Minute()))

	for table := range c.tablesAndStartTime {
		c.updateTableUsage(table, c.endTime)
	}

	for table, duration := range c.tablesAndHours {
		price := c.calculatePrice(duration)
		c.outputs = append(c.outputs, fmt.Sprintf("%d %d %02d:%02d", table, price, int(duration.Hours()), int(duration.Minutes())%60))
	}
	c.outputs = append(c.outputs, c.outputsAfterEnd...)
}

// calculatePrice вычисляет стоимость использования стола за определенное время.
func (c *computerClub) calculatePrice(duration time.Duration) int {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	if minutes > 30 {
		return (hours + 2) * c.costOfHour
	}
	return (hours + 1) * c.costOfHour
}

// GetOutput возвращает выходные данные, полученные в результате обработки событий.
func (c *computerClub) GetOutput() []string {
	return c.outputs
}
