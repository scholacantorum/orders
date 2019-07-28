//
//  Store.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-23.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import Foundation

struct Allow {
    var card = false
    var cash = false
    var willcall = false
}

struct Event: Codable {
    var id: String
    var name: String
    var start: String
    var freeEntries: [String]?
}

struct Product: Codable {
    var id: String
    var name: String
    var message: String?
    var price: Int!
    var ticketCount: Int!
}

var store = Store()

class Store {

    var username = ""
    var auth = ""
    var allow = Allow()
    var baseURL = ""
    var event: Event!
    var products: [Product] = []
    var logoutCallback: () -> Void = {}
    var admitted = 0
    var sold = 0
    var cash = 0
    var check = 0

    func logout() {
        username = ""
        auth = ""
        allow.card = false
        allow.cash = false
        allow.willcall = false
        baseURL = ""
        event = nil
        products = []
        admitted = 0
        sold = 0
        cash = 0
        check = 0
        cardReaderWatcher.disconnect()
        logoutCallback()
    }
}
