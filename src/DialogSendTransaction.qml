import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12

Dialog {
    id: dialogSend
    anchors.centerIn: Overlay.overlay

    title: qsTr("Password requested")
    modal: true
    clip: true
    focus: true
    standardButtons: Dialog.Ok | Dialog.Cancel

    onOpened: {
        dialogSend.forceActiveFocus()
        textFieldPassword.forceActiveFocus()
        standardButton(Dialog.Ok).enabled = textFieldPassword.text !== ""
    }

    onClosed: {
        textFieldPassword.clear()
    }

    ColumnLayout {
        anchors.fill: parent

        TransactionDetails {
            id: transactionDetail

            // EXAMPLE
            date: "2019-02-26 15:27"
            status: TransactionDetails.Status.Preview
            type: TransactionDetails.Type.Receive
            amount: 10
            hoursReceived: 16957
            hoursBurned: 33901
            transactionID: "kq7wdkkUT736dnuyQasdhsaDJ9874jk"

            Layout.fillWidth: true
        }

        TextField {
            id: textFieldPassword

            placeholderText: qsTr("Enter the password")
            echoMode: TextField.Password
            Layout.fillWidth: true
            focus: true

            onTextChanged: {
                dialogSend.standardButton(Dialog.Ok).enabled = text !== ""
            }
        }
    }
}
