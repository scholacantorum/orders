//
//  Main.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-24.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

// For testing, to skip login page:
// import Stripe

class Main: UIViewController {

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
        // store.event = Event(id: "2019-07-29", name: "Summer Sing", start: "2019-07-29T19:30:00", freeEntries: ["Student"])
        // store.products = [
        //     Product(id: "ticket-2019-07-29", name: "July 29", message: nil, price: 1700, ticketCount: 1),
        //     Product(id: "summer-sings-2019", name: "Flex Pass", message: nil, price: 8500, ticketCount: 6),
        //     Product(id: "summer-sings-2019-student", name: "Student", message: nil, price: 0, ticketCount: 1),
        // ]
        // backend.noop()

        let tindex = store.event.start.firstIndex(of: "T")!
        navigationItem.title = "\(store.event.start[..<tindex]) \(store.event.name)"
        navigationItem.hidesBackButton = true
        view.backgroundColor = UIColor.white

        let menuController = MainMenu(reconnect: reconnect)
        let innerController = UINavigationController(rootViewController: menuController)
        innerController.setNavigationBarHidden(true, animated: false)
        addChild(innerController)
        view.addSubview(innerController.view)

        NSLayoutConstraint.useAndActivateConstraints([
            innerController.view.leftAnchor.constraint(equalTo: view.leftAnchor),
            innerController.view.rightAnchor.constraint(equalTo: view.rightAnchor),
            innerController.view.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
            innerController.view.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor),
        ])
    }

    func reconnect() {
        navigationController!.popViewController(animated: true)
    }

}
