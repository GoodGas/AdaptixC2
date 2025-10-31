#ifndef ADAPTIXCLIENT_AGENT_H
#define ADAPTIXCLIENT_AGENT_H

#include <main.h>

class Commander;
class ConsoleWidget;
class BrowserFilesWidget;
class BrowserProcessWidget;
class TerminalWidget;
class AdaptixWidget;
class AgentTableWidgetItem;
class GraphItem;

class Agent
{
public:
    AdaptixWidget* adaptixWidget = nullptr;

    AgentData data = {};

    QImage imageActive   = QImage();
    QImage imageInactive = QImage();
    QIcon  iconOs        = QIcon();

    QString LastMark   = QString();
    QString LastUpdate = QString();

    QString connType     = QString();
    QString listenerType = QString();
    QString parentId     = QString();
    QVector<QString> childsId;

    GraphItem* graphItem  = nullptr;
    QImage     graphImage = QImage();

    Commander*            commander      = nullptr;
    ConsoleWidget*        Console        = nullptr;
    BrowserFilesWidget*   FileBrowser    = nullptr;
    BrowserProcessWidget* ProcessBrowser = nullptr;
    TerminalWidget*       Terminal       = nullptr;

    bool active = true;
    bool show   = true;
    QColor bg_color = QColor();
    QColor fg_color = QColor();

    explicit Agent(QJsonObject jsonObjAgentData, AdaptixWidget* w );
    ~Agent();

    void    Update(QJsonObject jsonObjAgentData);
    void    MarkItem(const QString &mark);
    void    UpdateImage();
    QString TasksCancel(const QStringList &tasks) const;
    QString TasksDelete(const QStringList &tasks) const;

    void SetParent(const PivotData &pivotData);
    void UnsetParent(const PivotData &pivotData);
    void AddChild(const PivotData &pivotData);
    void RemoveChild(const PivotData &pivotData);
};

#endif
