package main

import (
    "testing"
    "io/ioutil"
)

func TestConfFileImport(t *testing.T) {
    ParseConfigFile("local.conf");
    if (Config.Rabbitmq_pass != "admin") {
        t.Errorf("Expected rabbitmq_pass to be admin but got %s", Config.Rabbitmq_pass)
    }

    if (Config.Rabbitmq_user != "admin") {
        t.Errorf("Expected rabbitmq_user to be admin but got %s", Config.Rabbitmq_user)
    }

    if (Config.Rabbitmq_port != 5672) {
        t.Errorf("Expected rabbitmq_port to be 5672 but got %d", Config.Rabbitmq_port)
    }

    if (Config.Rabbitmq_host != "10.0.0.142") {
        t.Errorf("Expected rabbitmq_host to be 10.0.0.142 but got %s", Config.Rabbitmq_host)
    }
}

func TestParseData(t *testing.T) {
    data, _ := ioutil.ReadFile("testdata.json")
    addStockData(data)

    if StockContainer.Data[0].Name != "F5 Networks Inc" {
        t.Errorf("Expected name to be F5 Networks Inc but got %s", StockContainer.Data[0].Name)
    }

    if StockContainer.Data[0].Symbol != "FFIV" {
        t.Errorf("Expected symbol to be FFIV but got %s", StockContainer.Data[0].Symbol)
    }

    if StockContainer.Data[0].Change != -3.13 {
        t.Errorf("Expected change to be -3.13 but got %d", StockContainer.Data[0].Change)
    }

    if StockContainer.Data[0].ChangePercent != -3.29404335929278 {
        t.Errorf("Expected change to be -3.29404335929278 but got %d", StockContainer.Data[0].ChangePercent)
    }

    if StockContainer.Data[0].Price != 91.89 {
        t.Errorf("Expected prive to be 91.89 but got %d", StockContainer.Data[0].Price)
    }
}