//
//  TicketQuantity.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-24.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

protocol TicketQuantityDelegate {
    func ticketQuantityChange(product: Product, sellQty: Int, useQty: Int) -> Void
}

let useButtonColor = UIColor(red: 23.0/255.0, green: 162.0/255.0, blue: 184.0/255.0, alpha: 1.0)

class TicketQuantity: UIViewController {

    var product: Product!
    var delegate: TicketQuantityDelegate!
    var sellqty = 0
    var useqty = 0
    var showUseRow = false

    lazy var sellRow = UIView()
    lazy var nameLabel: UILabel = {
        var view = UILabel()
        view.font = UIFont.systemFont(ofSize: 24.0)
        return view
    }()
    lazy var sellDownButton: UIButton = {
        let view = UIButton()
        view.setTitle("–", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 36.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(sellDownButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var sellQtyLabel: UILabel = {
        let view = UILabel()
        view.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.textAlignment = .center
        return view
    }()
    lazy var sellUpButton: UIButton = {
        let view = UIButton()
        view.setTitle("+", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 36.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = scholaBlue
        view.addTarget(self, action: #selector(sellUpButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var priceLabel: UILabel = {
        var view = UILabel()
        view.font = UIFont.systemFont(ofSize: 24.0)
        view.textAlignment = .right
        return view
    }()
    lazy var useRow = UIView()
    lazy var andUseLabel: UILabel = {
        var view = UILabel()
        view.text = "and use"
        view.font = UIFont.systemFont(ofSize: 24.0)
        return view
    }()
    lazy var useDownButton: UIButton = {
        let view = UIButton()
        view.setTitle("–", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 36.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = useButtonColor
        view.addTarget(self, action: #selector(useDownButton(_:)), for: .touchUpInside)
        return view
    }()
    lazy var useQtyLabel: UILabel = {
        let view = UILabel()
        view.font = UIFont.boldSystemFont(ofSize: 24.0)
        view.textAlignment = .center
        return view
    }()
    lazy var useUpButton: UIButton = {
        let view = UIButton()
        view.setTitle("+", for: .normal)
        view.titleLabel!.font = UIFont.boldSystemFont(ofSize: 36.0)
        view.setTitleColor(UIColor.white, for: .normal)
        view.layer.cornerRadius = 5.0
        view.backgroundColor = useButtonColor
        view.addTarget(self, action: #selector(useUpButton(_:)), for: .touchUpInside)
        return view
    }()

    var noUseBottomConstraint: NSLayoutConstraint!

    init(product: Product, delegate: TicketQuantityDelegate) {
        super.init(nibName: nil, bundle: nil)
        self.product = product
        self.delegate = delegate
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()
        view.addSubview(sellRow)
        sellRow.addSubview(nameLabel)
        sellRow.addSubview(sellDownButton)
        sellRow.addSubview(sellQtyLabel)
        sellRow.addSubview(sellUpButton)
        sellRow.addSubview(priceLabel)
        noUseBottomConstraint = sellRow.bottomAnchor.constraint(equalTo: view.bottomAnchor)
        NSLayoutConstraint.useAndActivateConstraints([
            sellRow.leftAnchor.constraint(equalTo: view.leftAnchor),
            sellRow.rightAnchor.constraint(equalTo: view.rightAnchor),
            sellRow.topAnchor.constraint(equalTo: view.topAnchor),
            noUseBottomConstraint,
            nameLabel.leftAnchor.constraint(equalTo: sellRow.leftAnchor, constant: 9.0),
            nameLabel.centerYAnchor.constraint(equalTo: sellRow.centerYAnchor),
            priceLabel.rightAnchor.constraint(equalTo: sellRow.rightAnchor, constant: -9.0),
            priceLabel.widthAnchor.constraint(equalToConstant: 60.0),
            priceLabel.centerYAnchor.constraint(equalTo: sellRow.centerYAnchor),
            sellUpButton.rightAnchor.constraint(equalTo: priceLabel.leftAnchor, constant: -9.0),
            sellUpButton.widthAnchor.constraint(equalToConstant: 45.0),
            sellUpButton.topAnchor.constraint(equalTo: sellRow.topAnchor),
            sellUpButton.bottomAnchor.constraint(equalTo: sellRow.bottomAnchor),
            sellUpButton.heightAnchor.constraint(equalToConstant: 45.0),
            sellQtyLabel.rightAnchor.constraint(equalTo: sellUpButton.leftAnchor, constant: -9.0),
            sellQtyLabel.widthAnchor.constraint(equalToConstant: 30.0),
            sellQtyLabel.centerYAnchor.constraint(equalTo: sellRow.centerYAnchor),
            sellDownButton.rightAnchor.constraint(equalTo: sellQtyLabel.leftAnchor, constant: -9.0),
            sellDownButton.widthAnchor.constraint(equalToConstant: 45.0),
            sellDownButton.topAnchor.constraint(equalTo: sellRow.topAnchor),
            sellDownButton.bottomAnchor.constraint(equalTo: sellRow.bottomAnchor),
            sellDownButton.heightAnchor.constraint(equalToConstant: 45.0),
        ])
        nameLabel.text = product.name
        sellQtyLabel.text = "0"
        priceLabel.text = "$\(product.price! / 100)"
        priceLabel.textColor = UIColor.darkGray
    }

    func addUseRow() {
        if showUseRow {
            return
        }
        noUseBottomConstraint.isActive = false
        view.addSubview(useRow)
        useRow.addSubview(andUseLabel)
        useRow.addSubview(useDownButton)
        useRow.addSubview(useQtyLabel)
        useRow.addSubview(useUpButton)
        NSLayoutConstraint.useAndActivateConstraints([
            useRow.leftAnchor.constraint(equalTo: view.leftAnchor),
            useRow.rightAnchor.constraint(equalTo: view.rightAnchor),
            useRow.topAnchor.constraint(equalTo: sellRow.bottomAnchor, constant: 18.0),
            useRow.bottomAnchor.constraint(equalTo: view.bottomAnchor),
            useUpButton.rightAnchor.constraint(equalTo: sellRow.rightAnchor, constant: -78.0),
            useUpButton.widthAnchor.constraint(equalToConstant: 45.0),
            useUpButton.topAnchor.constraint(equalTo: useRow.topAnchor),
            useUpButton.bottomAnchor.constraint(equalTo: useRow.bottomAnchor),
            useUpButton.heightAnchor.constraint(equalToConstant: 45.0),
            useQtyLabel.rightAnchor.constraint(equalTo: useUpButton.leftAnchor, constant: -9.0),
            useQtyLabel.widthAnchor.constraint(equalToConstant: 30.0),
            useQtyLabel.centerYAnchor.constraint(equalTo: useRow.centerYAnchor),
            useDownButton.rightAnchor.constraint(equalTo: useQtyLabel.leftAnchor, constant: -9.0),
            useDownButton.widthAnchor.constraint(equalToConstant: 45.0),
            useDownButton.topAnchor.constraint(equalTo: useRow.topAnchor),
            useDownButton.bottomAnchor.constraint(equalTo: useRow.bottomAnchor),
            useDownButton.heightAnchor.constraint(equalToConstant: 45.0),
            andUseLabel.rightAnchor.constraint(equalTo: useDownButton.leftAnchor, constant: -9.0),
            andUseLabel.centerYAnchor.constraint(equalTo: useRow.centerYAnchor),
        ])
        showUseRow = true
    }

    @objc func sellDownButton(_ sender: UIButton) {
        if sellqty == 0 {
            return
        }
        sellqty -= 1
        if useqty > 0 {
            useqty -= 1
        }
        if sellqty > 0 {
            priceLabel.text = "$\(product.price! * sellqty / 100)"
            priceLabel.textColor = UIColor.black
        } else {
            priceLabel.text = "$\(product.price! / 100)"
            priceLabel.textColor = UIColor.darkGray
        }
        sellQtyLabel.text = "\(sellqty)"
        if showUseRow {
            useQtyLabel.text = "\(useqty)"
        }
        delegate.ticketQuantityChange(product: product, sellQty: sellqty, useQty: useqty)
    }

    @objc func sellUpButton(_ sender: UIButton) {
        sellqty += 1
        useqty += 1
        priceLabel.text = "$\(product.price! * sellqty / 100)"
        priceLabel.textColor = UIColor.black
        sellQtyLabel.text = "\(sellqty)"
        if sellqty > 1 || (sellqty == 1 && product.ticketCount > 1) {
            addUseRow()
        }
        if showUseRow {
            useQtyLabel.text = "\(useqty)"
        }
        delegate.ticketQuantityChange(product: product, sellQty: sellqty, useQty: useqty)
    }

    @objc func useDownButton(_ sender: UIButton) {
        if useqty == 0 {
            return
        }
        useqty -= 1
        useQtyLabel.text = "\(useqty)"
        delegate.ticketQuantityChange(product: product, sellQty: sellqty, useQty: useqty)
    }

    @objc func useUpButton(_ sender: UIButton) {
        if useqty >= product.ticketCount * sellqty {
            return
        }
        useqty += 1
        useQtyLabel.text = "\(useqty)"
        delegate.ticketQuantityChange(product: product, sellQty: sellqty, useQty: useqty)
    }
}
