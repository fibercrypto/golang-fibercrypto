import QtQuick 2.12
import QtQuick.Controls 2.12
import QtQuick.Controls.Material 2.12
import QtQuick.Layouts 1.12

// Resource imports
// import "qrc:/ui/src/ui/Controls"
// import "qrc:/ui/src/ui"
import "../" // For quick UI development, switch back to resources when making a release
import "../Controls" // For quick UI development, switch back to resources when making a release

Item {
    id: addressListDelegate

    signal qrCodeRequested(var data)
    signal addressTextChanged(string text)
    signal numberOfAddressesChanged(int count)

    onQrCodeRequested: {
        dialogQR.setVars(data)
        dialogQR.open()
    }

    implicitHeight: rootLayout.height
    clip: true

    RowLayout {
        id: rootLayout
        width: addressListDelegate.width
        clip: true
        spacing: 20
        opacity: 0.0

        Component.onCompleted: { opacity = 1.0 } // Not the best way to do this
        Behavior on opacity { NumberAnimation { duration: 500; easing.type: Easing.OutQuint } }

        RowLayout {
            Layout.fillWidth: true
            spacing: 8

            ToolButtonQR {
                id: toolButtonQR

                iconSize: "24x24"

                onClicked: {
                    qrCodeRequested(value)
                }
            }

            TextField {
                id: tfAddr
                font.family: "Code New Roman"
                placeholderText: qsTr("Address No.") + " " + (index + 1)
                text: value ? value : ""
                selectByMouse: true
                Layout.fillWidth: true
                Material.accent: abm.addressIsValid(text) ? parent.Material.accent : Material.color(Material.Red)
                onTextChanged: {
                    value = text
                    addressTextChanged(text)
                }
            }
        } // RowLayout

        RowLayout {
            ComboBox {
                id: cbCoinTypes
                model: coins
                currentIndex: getIndexForCoinType(coinType)
                onCurrentTextChanged: {
                    coinType = cbCoinTypes.currentText
                }
            }
        }

        ToolButton {
            id: toolButtonAddRemoveDestination
            // The 'accent' attribute is used for button highlighting
            Material.accent: index === 0 ? Material.Teal : Material.Red
            icon.source: "qrc:/images/resources/images/icons/" + (index === 0 ? "add" : "remove") + "-circle.svg"
            highlighted: true // enable the usage of the `Material.accent` attribute

            Layout.alignment: Qt.AlignRight

            onClicked: {
                if (index === 0) {
                    listModelAddresses.append( { "address": "", "coinType": "" } )
                } else {
                    listModelAddresses.remove(index)
                }
                numberOfAddressesChanged(listModelAddresses.count)
            }
        } // ToolButton
    } // RowLayout (rootLayout)

    ListModel {
        id:coins
        ListElement{
            type:"SKY"
        }
        // ListElement{
        //     type:"BTC"
        // }
    }

    function getIndexForCoinType(coinType) {
        var index = 0;
        for(var i = 0; i < coins.count; i++) {
            if (coins.get(i).type === coinType) {
                index = i
            }
        }
        return index;
    }
}
