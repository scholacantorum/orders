//
//  EventChooserController.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-24.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

// For testing, to skip login page:
// import Stripe

class EventChooser: UITableViewController {

    var events: [Event] = []

    override func viewDidLoad() {
        super.viewDidLoad()

        // For testing, skip login page:
        // store.username = "sroth"
        // store.auth = "9999-9999-9999"
        // store.baseURL = "https://orders-test.scholacantorum.org/api"
        // store.allow.card = true
        // store.allow.cash = true
        // store.allow.willcall = true
        // STPPaymentConfiguration.shared().publishableKey = "pk_test_QPwvhWbGaakWn7DGcco8J5Nd"

        navigationItem.title = "Choose Event"
        navigationItem.hidesBackButton = true
        tableView.register(UITableViewCell.self, forCellReuseIdentifier: "eventCell")
        backend.eventlist {events, error in
            DispatchQueue.main.async {
                if let error = error {
                    let alert = UIAlertController(title: "Server Error", message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    return
                }
                self.events = events!
                self.tableView.reloadData()
            }
        }
    }

    override func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return events.count
    }

    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "eventCell", for: indexPath)
        let event = events[indexPath.row]
        let firstSpace = event.start.firstIndex(of: "T")!
        let date = event.start[..<firstSpace]
        cell.textLabel!.text = date + " " + event.name
        return cell
    }

    override func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        store.event = events[indexPath.row]
        backend.eventProducts {products, error in
            DispatchQueue.main.async {
                if let error = error {
                    let alert = UIAlertController(title: "Server Error", message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    return
                }
                let filteredProducts = products!.filter { product in
                    return product.message == nil
                }
                if filteredProducts.isEmpty {
                    let alert = UIAlertController(title: "Not for Sale", message: "No tickets are on sale for that event.", preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    return
                }
                store.products = filteredProducts
                if store.allow.card {
                    self.navigationController!.pushViewController(ChooseCardReader(), animated: true)
                } else {
                    self.navigationController!.pushViewController(Main(), animated: true)
                }
            }
        }
    }

}
