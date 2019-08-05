//
//  Backend.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-23.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import Foundation
import Stripe
import StripeTerminal

let backend = Backend()

struct Order: Codable {
    var id: Int?
    var source: String
    var name: String?
    var email: String?
    var payments: [OrderPayment]
    var lines: [OrderLine]
    var error: String?
}
struct OrderPayment: Codable {
    var type: String
    var subtype: String?
    var method: String?
    var amount: Int
}
struct OrderLine: Codable {
    var product: String
    var quantity: Int
    var used: Int?
    var usedAt: String?
    var price: Int
}
struct WillCallOrder: Codable {
    var id: Int
    var name: String
}
struct TicketUsage: Codable {
    var id: Int
    var name: String?
    var error: String?
    var scan: String
    var classes: [TicketClassUsage]
}
struct TicketClassUsage: Codable {
    var name: String
    var min: Int
    var max: Int
    var used: Int
    var overflow: Bool?
}

class Backend: ConnectionTokenProvider {

    init() {
        Terminal.setTokenProvider(self)
    }

    struct LoginResult: Codable {
        var privScanTickets: Bool
        var privInPersonSales: Bool
        var privViewOrders: Bool
        var stripePublicKey: String
        var token: String
    }

    func login(_ username: String, _ password: String, _ testmode: Bool, _ allow: Allow, handler: @escaping (String?) -> Void) {
        let baseurl = "https://orders\(testmode ? "-test" : "").scholacantorum.org/api"
        let url = URL(string: baseurl + "/login")
        var request = URLRequest(url: url!)
        request.httpMethod = "POST"
        request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")
        let usernameEncoded = username.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!
        let passwordEncoded = password.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!
        request.httpBody = "username=\(usernameEncoded)&password=\(passwordEncoded)".data(using: String.Encoding.utf8)
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error with login request: \(error)")
                handler("Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error with login request: no response")
                handler("Server error: no response")
                return
            }
            if response.statusCode == 401 {
                handler("Login incorrect")
                return
            }
            if response.statusCode != 200 {
                print("Error with login request: \(response.statusCode)")
                handler("Server error: \(response.statusCode)")
                return
            }
            var result: LoginResult
            do {
                result = try JSONDecoder().decode(LoginResult.self, from: data!)
            } catch {
                print("Error decoding response to login request: \(error)")
                handler("Server error: \(error)")
                return
            }
            if !result.privScanTickets {
                handler("Not authorized to use this app")
                return
            }
            if (allow.card || allow.cash) && !result.privInPersonSales {
                handler("Not authorized to sell tickets")
                return
            }
            if allow.willcall && !result.privViewOrders {
                handler("Not authorized to view will call list")
                return
            }
            store.username = username
            store.auth = result.token
            store.allow = allow
            store.baseURL = baseurl
            STPPaymentConfiguration.shared().publishableKey = result.stripePublicKey
            handler(nil)
        }
        task.resume()
    }

    func eventlist(_ handler: @escaping ([Event]?, String?) -> Void) {
        let url = URL(string: store.baseURL + "/event?future=1&freeEntries=1")
        var request = URLRequest(url: url!)
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error getting event list: \(error)")
                handler(nil, "Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error getting event list: no response")
                handler(nil, "Server error: no response")
                return
            }
            if response.statusCode != 200 {
                print("Error getting event list: \(response.statusCode)")
                handler(nil, "Server error: \(response.statusCode)")
                return
            }
            var result: [Event]
            do {
                result = try JSONDecoder().decode([Event].self, from: data!)
            } catch {
                print("Error decoding event list: \(error)")
                handler(nil, "Server error: \(error)")
                return
            }
            handler(result, nil)
        }
        task.resume()
    }

    struct eventProductsResponse: Codable {
        var coupon: Bool
        var products: [Product]
    }

    func eventProducts(_ handler: @escaping ([Product]?, String?) -> Void) {
        let url = URL(string: store.baseURL + "/prices?event=" + store.event.id.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!)
        var request = URLRequest(url: url!)
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error getting event products: \(error)")
                handler(nil, "Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error getting event products: no response")
                handler(nil, "Server error: no response")
                return
            }
            if response.statusCode != 200 {
                print("Error getting event products: \(response.statusCode)")
                handler(nil, "Server error: \(response.statusCode)")
                return
            }
            var result: eventProductsResponse
            do {
                result = try JSONDecoder().decode(eventProductsResponse.self, from: data!)
            } catch {
                print("Error decoding event products: \(error)")
                handler(nil, "Server error: \(error)")
                return
            }
            handler(result.products, nil)
        }
        task.resume()
    }

    func placeOrder(order: Order, handler: @escaping (Order?, String?) -> Void) {
        let url = URL(string: store.baseURL + "/order")
        var request = URLRequest(url: url!)
        request.httpMethod = "POST"
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        request.httpBody = try! JSONEncoder().encode(order)
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error placing order: \(error)")
                handler(nil, "Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error placing order: no response")
                handler(nil, "Server error: no response")
                return
            }
            if response.statusCode == 400 {
                handler(nil, String(decoding: data!, as: UTF8.self))
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 200 {
                print("Error placing order: \(response.statusCode)")
                handler(nil, "Server error: \(response.statusCode)")
                return
            }
            var result: Order
            do {
                result = try JSONDecoder().decode(Order.self, from: data!)
            } catch {
                print("Error decoding placed order: \(error)")
                handler(nil, "Server error: \(error)")
                return
            }
            if let error = result.error {
                handler(nil, error)
                return
            }
            handler(result, nil)
        }
        task.resume()
    }

    func cancelOrder(order: Order) {
        let url = URL(string: store.baseURL + "/order/\(order.id!)")
        var request = URLRequest(url: url!)
        request.httpMethod = "DELETE"
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error canceling order: \(error)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error canceling order: no response")
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 204 {
                print("Error canceling order: \(response.statusCode)")
            }
        }
        task.resume()
    }

    func captureOrderPayment(order: Order, handler: @escaping (Order?, String?) -> Void) {
        let url = URL(string: store.baseURL + "/order/\(order.id!)/capturePayment")
        var request = URLRequest(url: url!)
        request.httpMethod = "POST"
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error capturing payment: \(error)")
                handler(nil, "Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error capturing payment: no response")
                handler(nil, "Server error: no response")
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 200 {
                print("Error capturing payment: \(response.statusCode)")
                handler(nil, "Server error: \(response.statusCode)")
                return
            }
            var result: Order
            do {
                result = try JSONDecoder().decode(Order.self, from: data!)
            } catch {
                print("Error decoding captured payment order: \(error)")
                handler(nil, "Server error: \(error)")
                return
            }
            handler(result, nil)
        }
        task.resume()
    }

    func sendEmailReceipt(order: Order, email: String, handler: @escaping (String?) -> Void) {
        let url = URL(string: store.baseURL + "/order/\(order.id!)/sendReceipt?email=" + email.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!)
        var request = URLRequest(url: url!)
        request.httpMethod = "POST"
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error sending receipt: \(error)")
                handler(error.localizedDescription)
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error sending receipt: no response")
                handler("no response")
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 204 {
                print("Error sending receipt: \(response.statusCode)")
                handler("Error \(response.statusCode)")
                return
            }
            handler(nil)
        }
        task.resume()
    }

    func fetchConnectionToken(_ completion: @escaping ConnectionTokenCompletionBlock) {
        let url = URL(string: store.baseURL + "/stripe/connectTerminal")!
        var request = URLRequest(url: url)
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error getting connection token: \(error)")
                completion(nil, error)
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error getting connection token: no response")
                completion(nil, NSError(domain: "org.scholacantorum.orders.ScholaPOS", code: 2000, userInfo: nil))
                return
            }
            if response.statusCode == 401 {
                print("Error getting connection token: not logged in")
                completion(nil, NSError(domain: "org.scholacantorum.orders.ScholaPOS", code: 2000, userInfo: nil))
                store.logout()
                return
            }
            if response.statusCode != 200 {
                print("Error getting connection token: \(response.statusCode)")
                completion(nil, NSError(domain: "org.scholacantorum.orders.ScholaPOS", code: 2000, userInfo: nil))
                return
            }
            var token: String
            do {
                token = try JSONSerialization.jsonObject(with: data!, options: JSONSerialization.ReadingOptions.allowFragments) as! String
            } catch {
                print("Error decoding connection token: \(error)")
                completion(nil, NSError(domain: "org.scholacantorum.orders.ScholaPOS", code: 2000, userInfo: nil))
                return
            }
            completion(token, nil)
        }
        task.resume()
    }

    func willCallList(_ event: Event, handler: @escaping ([WillCallOrder]?, String?) -> Void) {
        let url = URL(string: store.baseURL + "/event/\(event.id.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!)/orders")
        var request = URLRequest(url: url!)
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error getting will call list: \(error)")
                handler(nil, "Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error getting will call list: no response")
                handler(nil, "Server error: no response")
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 200 {
                print("Error getting will call list: \(response.statusCode)")
                handler(nil, "Server error: \(response.statusCode)")
                return
            }
            var result: [WillCallOrder]
            do {
                result = try JSONDecoder().decode([WillCallOrder].self, from: data!)
            } catch {
                print("Error decoding will call list: \(error)")
                handler(nil, "Server error: \(error)")
                return
            }
            handler(result, nil)
        }
        task.resume()
    }

    func fetchTicketUsage(event: Event, tokenOrID: String, handler: @escaping (TicketUsage?, String?) -> Void) {
        let url = URL(string: store.baseURL + "/event/\(event.id.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!)/ticket/\(tokenOrID)")
        var request = URLRequest(url: url!)
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error getting ticket usage: \(error)")
                handler(nil, "Server error: \(error.localizedDescription)")
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error getting ticket usage: no response")
                handler(nil, "Server error: no response")
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 200 {
                print("Error getting ticket usage: \(response.statusCode)")
                handler(nil, "Server error: \(response.statusCode)")
                return
            }
            var result: TicketUsage
            do {
                result = try JSONDecoder().decode(TicketUsage.self, from: data!)
            } catch {
                print("Error decoding ticket usage: \(error)")
                handler(nil, "Server error: \(error)")
                return
            }
            if let error = result.error {
                handler(nil, error)
                return
            }
            handler(result, nil)
        }
        task.resume()
    }

    func useTickets(event: Event, usage: TicketUsage, handler: @escaping (String?) -> Void) {
        let url = URL(string: store.baseURL + "/event/\(event.id.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!)/ticket/\(usage.id)")
        var request = URLRequest(url: url!)
        request.httpMethod = "POST"
        request.setValue(store.auth, forHTTPHeaderField: "Auth")
        request.setValue("application/x-www-form-urlencoded", forHTTPHeaderField: "Content-Type")
        var body = "scan=" + usage.scan
        for cl in usage.classes {
            body += "&class=\(cl.name.addingPercentEncoding(withAllowedCharacters: NSCharacterSet.urlQueryAllowed)!)&used=\(cl.used)"
        }
        request.httpBody = body.data(using: String.Encoding.utf8)
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Error using tickets: \(error)")
                handler(error.localizedDescription)
                return
            }
            guard let response = response as? HTTPURLResponse else {
                print("Error using tickets: no response")
                handler("no response")
                return
            }
            if response.statusCode == 401 {
                store.logout()
                return
            }
            if response.statusCode != 200 {
                print("Error using tickets: \(response.statusCode)")
                handler("Error \(response.statusCode)")
                return
            }
            handler(nil)
        }
        task.resume()
    }

    func noop() {
    }

}
