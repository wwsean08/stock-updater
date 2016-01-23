package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"

//    "stocks/vendor/github.com/spf13/cobra"
//    "stocks/vendor/github.com/streadway/amqp"
    "github.com/spf13/cobra"
    "github.com/streadway/amqp"
)

var cfgFile string
var baseEndpoint = "http://dev.markitondemand.com/MODApis/Api/v2/Quote/json?symbol="

// Structure containing the config data
type ConfigData struct {
    Stocks        []string    `json:"stocks"`
    Rabbitmq_host string      `json:"rabbitmq_host"`
    Rabbitmq_port int         `json:"rabbitmq_port"`
    Rabbitmq_user string      `json:"rabbitmq_user"`
    Rabbitmq_pass string      `json:"rabbitmq_pass"`
}

type StockData struct {
    Name          string  `json:"name"`
    Symbol        string  `json:"symbol"`
    Price         float32 `json:"price"`
    Change        float32 `json:"change"`
    ChangePercent float32 `json:"changePercent"`
}

type StockDataContainer struct {
    Data []StockData `json:"stockData"`
}

var Config ConfigData;
var StockContainer = StockDataContainer{}

func main() {
    var cmdMain = &cobra.Command{
        Use: "stock",
        Short: "Get stock prices of stocks and send them to rabbitmq",
        Long: "Get the stock prices of a list of given stocks and send them to rabitmq",
        Run:func(cmd *cobra.Command, args []string) {
            reportStocks();
        },
    }
    cmdMain.PersistentFlags().StringVar(&cfgFile, "config", "", "config file with stock and connection information")
    cmdMain.Execute()
}

func reportStocks() {
    ParseConfigFile(cfgFile)
    for _, symbol := range Config.Stocks {
        var endpoint = baseEndpoint + symbol
        resp, err := http.Get(endpoint)
        failOnError(err, "Error performing get request to " + endpoint)
        content, err := ioutil.ReadAll(resp.Body)
        failOnError(err, "Error occured reading body");
        resp.Body.Close()
        addStockData(content)
    }
    message, err := json.Marshal(StockContainer)
    failOnError(err, "Failed to create json message")
    sendMessageToRabbitMQ(message)
}

func addStockData(content []byte) {
    type MODData struct {
        Name          string `json:"Name"`
        Symbol        string `json:"Symbol"`
        Price         float32 `json:"LastPrice"`
        Change        float32 `json:"Change"`
        ChangePercent float32 `json:"ChangePercent"`
    }
    var moddata = MODData{}
    failOnError(json.Unmarshal(content, &moddata), "Error parsing data from Markit On Demand")
    var stock = StockData{
        Name:           moddata.Name,
        Symbol:         moddata.Symbol,
        Change:         moddata.Change,
        ChangePercent:  moddata.ChangePercent,
        Price:          moddata.Price,
    }
    StockContainer.Data = append(StockContainer.Data, stock)
}

// Connects to the RabbitMQ host and sends a message containing the stock
// information for the stock symbols in the config file.
func sendMessageToRabbitMQ(message []byte) {
    var amqpAddress = "amqp://" + Config.Rabbitmq_user + ":" + Config.Rabbitmq_pass + "@" + Config.Rabbitmq_host + ":" + strconv.Itoa(Config.Rabbitmq_port)
    conn, err := amqp.Dial(amqpAddress)
    failOnError(err, "Failed to connect to RabbitMQ")

    ch, err := conn.Channel()
    failOnError(err, "Failed to open channel")
    defer ch.Close();

    q, err := ch.QueueDeclare(
        "stock", // name
        false, // durable
        true, // delete when unused
        false, // exclusive
        false, // no-wait
        nil, // arguments
    )
    failOnError(err, "Failed to declare queue")

    err = ch.Publish(
        "amq.topic",
        q.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body: (message),
        })
    failOnError(err, "Failed to publish stock update")
}

// Parse the config file from filePath and create a Config struct
func ParseConfigFile(filePath string) {
    content, err := ioutil.ReadFile(filePath)
    failOnError(err, "Unable to open file")
    failOnError(json.Unmarshal(content, &Config), "Unable to parse json data")
}

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}