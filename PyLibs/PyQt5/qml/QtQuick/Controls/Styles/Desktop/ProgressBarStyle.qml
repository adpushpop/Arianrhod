/****************************************************************************
**
** Copyright (C) 2015 The Qt Company Ltd.
** Contact: http://www.qt.io/licensing/
**
** This file is part of the Qt Quick Controls module of the Qt Toolkit.
**
** $QT_BEGIN_LICENSE:LGPL3$
** Commercial License Usage
** Licensees holding valid commercial Qt licenses may use this file in
** accordance with the commercial license agreement provided with the
** Software or, alternatively, in accordance with the terms contained in
** a written agreement between you and The Qt Company. For licensing terms
** and conditions see http://www.qt.io/terms-conditions. For further
** information use the contact form at http://www.qt.io/contact-us.
**
** GNU Lesser General Public License Usage
** Alternatively, this file may be used under the terms of the GNU Lesser
** General Public License version 3 as published by the Free Software
** Foundation and appearing in the file LICENSE.LGPLv3 included in the
** packaging of this file. Please review the following information to
** ensure the GNU Lesser General Public License version 3 requirements
** will be met: https://www.gnu.org/licenses/lgpl.html.
**
** GNU General Public License Usage
** Alternatively, this file may be used under the terms of the GNU
** General Public License version 2.0 or later as published by the Free
** Software Foundation and appearing in the file LICENSE.GPL included in
** the packaging of this file. Please review the following information to
** ensure the GNU General Public License version 2.0 requirements will be
** met: http://www.gnu.org/licenses/gpl-2.0.html.
**
** $QT_END_LICENSE$
**
****************************************************************************/

import QtQuick 2.2
import QtQuick.Controls 1.2
import QtQuick.Controls.Private 1.0

Style {
    property Component panel: StyleItem {
        anchors.fill: parent
        elementType: "progressbar"
        // XXX: since desktop uses int instead of real, the progressbar
        // range [0..1] must be stretched to a good precision
        property int factor : 1000
        property int decimals: 3
        value:   indeterminate ? 0 : control.value.toFixed(decimals) * factor // does indeterminate value need to be 1 on windows?
        minimum: indeterminate ? 0 : control.minimumValue.toFixed(decimals) * factor
        maximum: indeterminate ? 0 : control.maximumValue.toFixed(decimals) * factor
        enabled: control.enabled
        horizontal: control.orientation === Qt.Horizontal
        hints: control.styleHints
        contentWidth: horizontal ? 200 : 23
        contentHeight: horizontal ? 23 : 200
    }
}
