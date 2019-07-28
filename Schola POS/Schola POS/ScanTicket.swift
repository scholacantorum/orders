//
//  ScanTicket.swift
//  Schola POS
//
//  Created by Steven Roth on 2019-07-27.
//  Copyright © 2019 Schola Cantorum. All rights reserved.
//

import UIKit
import AVFoundation

let qrcodeRE = ".*/ticket/\\d\\d\\d\\d-\\d\\d\\d\\d-\\d\\d\\d\\d$"
let qrcodePred = NSPredicate(format: "SELF MATCHES %@", qrcodeRE)

class ScanTicket: UIViewController, AVCaptureMetadataOutputObjectsDelegate {

    var captureSession: AVCaptureSession!
    var previewLayer: AVCaptureVideoPreviewLayer!

    override func viewDidLoad() {
        super.viewDidLoad()
        view.backgroundColor = UIColor.white

        captureSession = AVCaptureSession()
        guard let videoCaptureDevice = AVCaptureDevice.default(for: .video) else {
            failed()
            return
        }
        let videoInput: AVCaptureDeviceInput
        do {
            videoInput = try AVCaptureDeviceInput(device: videoCaptureDevice)
        } catch {
            failed()
            return
        }
        if captureSession.canAddInput(videoInput) {
            captureSession.addInput(videoInput)
        } else {
            failed()
            return
        }
        let metadataOutput = AVCaptureMetadataOutput()
        if captureSession.canAddOutput(metadataOutput) {
            captureSession.addOutput(metadataOutput)
            metadataOutput.setMetadataObjectsDelegate(self, queue: DispatchQueue.main)
            metadataOutput.metadataObjectTypes = [.qr]
        } else {
            failed()
            return
        }
        previewLayer = AVCaptureVideoPreviewLayer(session: captureSession)
        previewLayer.frame = CGRect(x: view.layer.bounds.minX, y: view.layer.bounds.minY, width: view.layer.bounds.width, height: view.layer.bounds.height - 100.0)
        previewLayer.videoGravity = .resizeAspectFill
        view.layer.addSublayer(previewLayer)

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

        NSLayoutConstraint.useAndActivateConstraints([
            cancelButton.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            cancelButton.widthAnchor.constraint(equalToConstant: 150.0),
            cancelButton.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor, constant: -9.0),
        ])

        captureSession.startRunning()
    }

    func failed() {
        let alert = UIAlertController(title: "Camera Error", message: "Unable to use camera.", preferredStyle: .alert)
        alert.addAction(UIAlertAction(title: "OK", style: .default))
        self.present(alert, animated: true, completion: nil)
        navigationController!.popViewController(animated: false)
    }

    @objc func cancelButton(_ sender: UIButton) {
        navigationController!.popViewController(animated: true)
    }

    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)
        if captureSession?.isRunning == false {
            captureSession.startRunning()
        }
    }

    override func viewWillDisappear(_ animated: Bool) {
        super.viewWillDisappear(animated)
        if captureSession?.isRunning == true {
            captureSession.stopRunning()
        }
    }

    func metadataOutput(_ output: AVCaptureMetadataOutput, didOutput metadataObjects: [AVMetadataObject], from connection: AVCaptureConnection) {
        if let metadataObject = metadataObjects.first {
            guard let readableObject = metadataObject as? AVMetadataMachineReadableCodeObject else { return }
            guard let stringValue = readableObject.stringValue else { return }
            if !qrcodePred.evaluate(with: stringValue) {
                let alert = UIAlertController(title: "Camera Error", message: "This barcode is not from a Schola order.", preferredStyle: .alert)
                alert.addAction(UIAlertAction(title: "OK", style: .default))
                self.present(alert, animated: true, completion: nil)
                return
            }
            captureSession.stopRunning()
            let token = String(stringValue[stringValue.index(stringValue.endIndex, offsetBy: -14)...])
            navigationController!.pushViewController(UseOrder(tokenOrID: token), animated: true)
        }
    }

}
