//
//  SalesSummary.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-25.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class SalesSummary: UIViewController {

    var order: Order!

    init(order: Order) {
        super.init(nibName: nil, bundle: nil)
        self.order = order
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        var constraints: [NSLayoutConstraint] = []
        let hbox = UIView()
        let totalLabel = UILabel()
        totalLabel.text = "TOTAL   $\(order.payments[0].amount/100)"
        totalLabel.font = UIFont.systemFont(ofSize: 24.0)
        hbox.addSubview(totalLabel)
        let vbox = UIView()
        var topAnchor = vbox.topAnchor
        for line in order.lines {
            let label = UILabel()
            for prod in store.products {
                if prod.id == line.product {
                    label.text = "\(line.quantity)× \(prod.name)"
                    break
                }
            }
            if line.product == "donation" {
                label.text = "Donation"
            }
            vbox.addSubview(label)
            constraints.append(contentsOf: [
                label.leftAnchor.constraint(equalTo: vbox.leftAnchor, constant: 9.0),
                label.rightAnchor.constraint(lessThanOrEqualTo: vbox.rightAnchor),
                label.topAnchor.constraint(equalTo: topAnchor),
            ])
            topAnchor = label.bottomAnchor
        }
        hbox.addSubview(vbox)
        view.addSubview(hbox)
        constraints.append(contentsOf: [
            topAnchor.constraint(equalTo: vbox.bottomAnchor),
            vbox.leftAnchor.constraint(equalTo: hbox.leftAnchor),
            {
                let c = vbox.widthAnchor.constraint(equalToConstant: 1.0)
                c.priority = UILayoutPriority.defaultLow
                return c
            }(),
            vbox.bottomAnchor.constraint(equalTo: hbox.bottomAnchor),
            totalLabel.rightAnchor.constraint(equalTo: hbox.rightAnchor, constant: -9.0),
            totalLabel.bottomAnchor.constraint(equalTo: hbox.bottomAnchor),
            hbox.leftAnchor.constraint(equalTo: view.leftAnchor),
            hbox.rightAnchor.constraint(equalTo: view.rightAnchor),
            hbox.topAnchor.constraint(equalTo: view.topAnchor),
            hbox.heightAnchor.constraint(greaterThanOrEqualTo: vbox.heightAnchor),
            hbox.heightAnchor.constraint(greaterThanOrEqualTo: totalLabel.heightAnchor),
            hbox.bottomAnchor.constraint(equalTo: view.bottomAnchor),
            {
                let c = hbox.heightAnchor.constraint(equalToConstant: 1.0)
                c.priority = UILayoutPriority.defaultLow
                return c
            }(),
        ])
        NSLayoutConstraint.useAndActivateConstraints(constraints)
    }

}
