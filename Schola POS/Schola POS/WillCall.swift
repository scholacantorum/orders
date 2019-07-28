//
//  WillCall.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-27.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit

class WillCallCell: UITableViewCell {
    override init(style: UITableViewCell.CellStyle, reuseIdentifier: String?) {
        super.init(style: .subtitle, reuseIdentifier: reuseIdentifier)
    }
    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }
}

class WillCall: UIViewController, UITableViewDataSource, UISearchBarDelegate, UITableViewDelegate {

    var orders: [WillCallOrder] = []
    var filteredOrders: [WillCallOrder] = []
    var searchBar: UISearchBar!
    var tableView: UITableView!

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white

        var constraints = [NSLayoutConstraint]()

        searchBar = UISearchBar()
        searchBar.delegate = self
        searchBar.placeholder = "Search"
        searchBar.showsBookmarkButton = false
        searchBar.setContentHuggingPriority(.required, for: .vertical)
        searchBar.setContentCompressionResistancePriority(.required, for: .vertical)
        view.addSubview(searchBar)
        constraints.append(contentsOf: [
            searchBar.leftAnchor.constraint(equalTo: view.leftAnchor),
            searchBar.rightAnchor.constraint(equalTo: view.rightAnchor),
            searchBar.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor),
        ])

        tableView = UITableView()
        tableView.register(WillCallCell.self, forCellReuseIdentifier: "willCallCell")
        tableView.delegate = self
        tableView.dataSource = self
        view.addSubview(tableView)
        constraints.append(contentsOf: [
            tableView.leftAnchor.constraint(equalTo: view.leftAnchor),
            tableView.rightAnchor.constraint(equalTo: view.rightAnchor),
            tableView.topAnchor.constraint(equalTo: searchBar.bottomAnchor),
        ])

        let cancelButton = UIButton()
        cancelButton.setTitle("Cancel", for: .normal)
        cancelButton.titleLabel!.font = UIFont.boldSystemFont(ofSize: 20.0)
        cancelButton.setTitleColor(UIColor.white, for: .normal)
        cancelButton.backgroundColor = UIColor.darkGray
        cancelButton.layer.cornerRadius = 5.0
        cancelButton.addTarget(self, action: #selector(cancelButton(_:)), for: .touchUpInside)
        cancelButton.setContentHuggingPriority(.required, for: .vertical)
        cancelButton.setContentCompressionResistancePriority(.required, for: .vertical)
        view.addSubview(cancelButton)
        constraints.append(contentsOf: [
            cancelButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            cancelButton.widthAnchor.constraint(equalToConstant: 150.0),
            cancelButton.topAnchor.constraint(equalTo: tableView.bottomAnchor, constant: 9.0),
            cancelButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])
        NSLayoutConstraint.useAndActivateConstraints(constraints)

        backend.willCallList(store.event) { orders, error in
            DispatchQueue.main.async {
                if let error = error {
                    let alert = UIAlertController(title: "Server Error", message: error, preferredStyle: .alert)
                    alert.addAction(UIAlertAction(title: "OK", style: .default))
                    self.present(alert, animated: true, completion: nil)
                    return
                }
                self.orders = orders!
                self.filterOrders()
            }
        }
    }

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return filteredOrders.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {
        let cell = tableView.dequeueReusableCell(withIdentifier: "willCallCell", for: indexPath)
        let order = filteredOrders[indexPath.row]
        cell.textLabel!.text = order.name
        cell.detailTextLabel!.text = "#\(order.id)"
        return cell
    }

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        let order = filteredOrders[indexPath.row]
        navigationController!.pushViewController(UseOrder(tokenOrID: "\(order.id)"), animated: true)
    }

    func searchBar(_ searchBar: UISearchBar, textDidChange searchText: String) {
        filterOrders()
    }

    func filterOrders() {
        filteredOrders.removeAll(keepingCapacity: false)
        let searchText = (searchBar.text ?? "").lowercased()
        for o in orders {
            if searchText == "" || o.name.lowercased().contains(searchText) || searchText == "\(o.id)" {
                filteredOrders.append(o)
            }
        }
        tableView.reloadData()
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }
    
}
