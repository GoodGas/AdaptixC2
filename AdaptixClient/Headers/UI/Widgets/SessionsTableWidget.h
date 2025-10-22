#ifndef ADAPTIXCLIENT_SESSIONSTABLEWIDGET_H
#define ADAPTIXCLIENT_SESSIONSTABLEWIDGET_H

#include <main.h>
#include <MainAdaptix.h>
#include <Utils/CustomElements.h>
#include <UI/Widgets/AbstractDock.h>
#include <UI/Widgets/AdaptixWidget.h>
#include <Agent/Agent.h>
#include <Client/Settings.h>

#include <QSortFilterProxyModel>

class Agent;
class AdaptixWidget;

enum SessionsColumns {
    SC_AgentID,
    SC_AgentType,
    SC_External,
    SC_Listener,
    SC_Internal,
    SC_Domain,
    SC_Computer,
    SC_User,
    SC_Os,
    SC_Process,
    SC_Pid,
    SC_Tid,
    SC_Tags,
    SC_Last,
    SC_Sleep,
    SC_ColumnCount
};



class AgentsFilterProxyModel : public QSortFilterProxyModel
{
Q_OBJECT
    AdaptixWidget* adaptixWidget = nullptr;
    bool    searchVisible  = false;
    bool    onlyActive     = false;
    QString filter1;
    QString filter2;
    QString filter3;

    static bool anyFieldContains(const QStringList &fields, const QString &pattern) {
        if (pattern.isEmpty()) return true;
        QRegularExpression re(pattern, QRegularExpression::CaseInsensitiveOption);
        for (const QString &f : fields) {
            if (f.contains(re))
                return true;
        }
        return false;
    }

public:
    explicit AgentsFilterProxyModel(AdaptixWidget* adaptix, QObject* parent = nullptr) : QSortFilterProxyModel(parent), adaptixWidget(adaptix) {
        setDynamicSortFilter(true);
        setFilterCaseSensitivity(Qt::CaseInsensitive);
    }

    void setSearchVisible(bool visible) {
        if (searchVisible == visible) return;
        searchVisible = visible;
        invalidateFilter();
    }
    void setOnlyActive(bool onlyActive) {
        if (this->onlyActive == onlyActive) return;
        this->onlyActive = onlyActive;
        invalidateFilter();
    }
    void setFilter1(const QString& text) {
        if (filter1 == text) return;
        filter1 = text;
        invalidateFilter();
    }
    void setFilter2(const QString& text) {
        if (filter2 == text) return;
        filter2 = text;
        invalidateFilter();
    }
    void setFilter3(const QString& text) {
        if (filter3 == text) return;
        filter3 = text;
        invalidateFilter();
    }
    void updateVisible() {
        invalidateFilter();
    }

protected:
    bool filterAcceptsRow(int row, const QModelIndex &parent) const override {
        if (!adaptixWidget || !sourceModel())
            return true;

        QModelIndex idxId = sourceModel()->index(row, 0, parent);
        if (!idxId.isValid())
            return true;

        QString agentId = sourceModel()->data(idxId, Qt::DisplayRole).toString();
        if (agentId.isEmpty())
            return true;

        const Agent* agent = adaptixWidget->AgentsMap.value(agentId, nullptr);
        if (!agent || !agent->show)
            return false;

        const AgentData &a = agent->data;

        if (!searchVisible)
            return true;

        if (onlyActive) {
            if (a.Mark == "Terminated" || a.Mark == "Inactive" || a.Mark == "Disconnect")
                return false;
        }

        QString username = a.Username;
        if (a.Elevated)
            username = "* " + username;
        if (!a.Impersonated.isEmpty())
            username += " [" + a.Impersonated + "]";

        QStringList searchable = { a.Id, a.Name, a.Listener, a.ExternalIP, a.InternalIP, a.Process, a.OsDesc, a.Domain, a.Computer, username, a.Tags };

        bool matches = true;

        if (!filter1.isEmpty() && !anyFieldContains(searchable, filter1))
            matches = false;
        if (!filter2.isEmpty() && !anyFieldContains(searchable, filter2))
            matches = false;
        if (!filter3.isEmpty() && !anyFieldContains(searchable, filter3))
            matches = false;

        return matches;
    }
};



class AgentsTableModel : public QAbstractTableModel
{
Q_OBJECT
    AdaptixWidget*   adaptixWidget;
    QVector<QString> agentsId;

public:
    explicit AgentsTableModel(AdaptixWidget* w, QObject* parent = nullptr) : QAbstractTableModel(parent), adaptixWidget(w) {}

    int rowCount(const QModelIndex&) const override {
        return agentsId.size();
    }

    int columnCount(const QModelIndex&) const override {
        return SC_ColumnCount;
    }

    QVariant data(const QModelIndex &index, const int role) const override {
        if (!index.isValid())
            return {};

        QString agentId = agentsId.at(index.row());
        Agent*  agent   = adaptixWidget->AgentsMap.value(agentId, nullptr);
        if (!agent)
            return {};

        AgentData d = agent->data;

        if (role == Qt::DisplayRole) {
            switch (index.column()) {
                case SC_AgentID:   return d.Id;
                case SC_AgentType: return d.Name;
                case SC_External:  return d.ExternalIP;
                case SC_Listener:  return d.Listener;
                case SC_Internal:  return d.InternalIP;
                case SC_Domain:    return d.Domain;
                case SC_Computer:  return d.Computer;
                case SC_User:
                {
                    QString username = d.Username;
                    if ( d.Elevated ) username = "* " + username;
                    if ( d.Impersonated != "" ) username += " [" + d.Impersonated + "]";
                    return username;
                }
                case SC_Os:        return d.OsDesc;
                case SC_Process:
                {
                    QString process = d.Process;
                    if ( !d.Arch.isEmpty() )
                        process += QString(" (%2)").arg(d.Arch);
                    return process;
                }
                case SC_Pid:       return d.Pid;
                case SC_Tid:       return d.Tid;
                case SC_Tags:      return d.Tags;
                case SC_Last:
                {
                    if ( d.Mark.isEmpty() || d.Mark == "No response" || d.Mark == "No worktime" ) {
                        return agent->LastMark;
                    }
                    return UnixTimestampGlobalToStringLocalSmall(d.LastTick);
                }
                case SC_Sleep:
                {
                    if ( d.Mark.isEmpty() ) {
                        if ( !d.Async ) {
                            if ( agent->connType == "internal" )
                                return QString::fromUtf8("\u221E  \u221E");
                            else
                                return QString::fromUtf8("\u27F6\u27F6\u27F6");
                        }
                        return QString("%1 (%2%)").arg( FormatSecToStr(d.Sleep) ).arg(d.Jitter);
                    }
                    return d.Mark;
                }
            }
        }

        if (role == Qt::TextAlignmentRole) {
            switch (index.column()) {
                case SC_AgentID:
                case SC_AgentType:
                case SC_External:
                case SC_Internal:
                case SC_Listener:
                case SC_Domain:
                case SC_Computer:
                case SC_User:
                case SC_Os:
                // case SC_Process:
                case SC_Pid:
                case SC_Tid:
                case SC_Last:
                case SC_Sleep:
                    return Qt::AlignCenter;
            }
        }

        if (role == Qt::DecorationRole) {
            if (index.column() == SC_Os) {
                return agent->iconOs;
            }
        }

        if (role == Qt::BackgroundRole) {
            if (agent->bg_color.isValid())
                return agent->bg_color;
            return QVariant();
        }
        if (role == Qt::ForegroundRole) {
            if (agent->fg_color.isValid())
                return agent->fg_color;
            return QVariant();
        }

        if (role == Qt::ToolTipRole) {
            if (index.column() == SC_Sleep) {
                QString WorkAndKill = "";
                if (d.WorkingTime || d.KillDate) {
                    if (d.WorkingTime) {
                        uint startH = ( d.WorkingTime >> 24 ) % 64;
                        uint startM = ( d.WorkingTime >> 16 ) % 64;
                        uint endH   = ( d.WorkingTime >>  8 ) % 64;
                        uint endM   = ( d.WorkingTime >>  0 ) % 64;

                        QChar c = QLatin1Char('0');
                        WorkAndKill = QString("Work time: %1:%2 - %3:%4\n").arg(startH, 2, 10, c).arg(startM, 2, 10, c).arg(endH, 2, 10, c).arg(endM, 2, 10, c);
                    }
                    if (d.KillDate) {
                        QDateTime dateTime = QDateTime::fromSecsSinceEpoch(d.KillDate);
                        WorkAndKill += QString("Kill date: %1").arg(dateTime.toString("dd.MM.yyyy hh:mm:ss"));
                    }
                }
                return WorkAndKill;
            }
        }

        return QVariant();
    }

    QVariant headerData(int section, Qt::Orientation o, int role) const override {
        if (role != Qt::DisplayRole || o != Qt::Horizontal)
            return {};

        static QStringList headers = {
            "Agent Id", "Type", "External", "Listener", "Internal",
            "Domain", "Computer", "User", "OS", "Process",
            "PID", "TID", "Tags", "Last", "Sleep"
        };

        return headers.value(section);
    }

    void add(const QString &agentId) {
        int rows = agentsId.size();

        beginInsertRows(QModelIndex(), rows, rows);
        agentsId.append(agentId);
        endInsertRows();
    }

    void update(const QString &agentId) {
        int row = agentsId.indexOf(agentId);
        if (row < 0)
            return;

        Q_EMIT dataChanged(index(row, 0), index(row, SC_ColumnCount - 1), { Qt::DisplayRole, Qt::ForegroundRole, Qt::BackgroundRole });
    }

    void remove(const QString &agentId) {
        int row = agentsId.indexOf(agentId);
        if (row < 0)
            return;

        beginRemoveRows(QModelIndex(), row, row);
        agentsId.removeAt(row);
        endRemoveRows();
    }

    void clear() {
        beginResetModel();
        agentsId.clear();
        endResetModel();
    }
};



class SessionsTableWidget : public DockTab
{
Q_OBJECT
    AdaptixWidget* adaptixWidget = nullptr;

    QGridLayout*  mainGridLayout = nullptr;
    QTableView*   tableView      = nullptr;
    QMenu*        menuSessions   = nullptr;
    QShortcut*    shortcutSearch = nullptr;

    AgentsTableModel*       agentsModel = nullptr;
    AgentsFilterProxyModel* proxyModel  = nullptr;

    QWidget*        searchWidget    = nullptr;
    QHBoxLayout*    searchLayout    = nullptr;
    QCheckBox*      checkOnlyActive = nullptr;
    QLineEdit*      inputFilter1    = nullptr;
    QLineEdit*      inputFilter2    = nullptr;
    QLineEdit*      inputFilter3    = nullptr;
    ClickableLabel* hideButton      = nullptr;

    void createUI();

public:
    QTimer* refreshTimer = nullptr;

    explicit SessionsTableWidget( AdaptixWidget* w );
    ~SessionsTableWidget() override;

    void AddAgentItem(Agent* newAgent) const;
    void UpdateAgentItem(const AgentData &oldDatam, const Agent* agent) const;
    void RemoveAgentItem(const QString &agentId) const;

    void UpdateColumnsVisible() const;
    void UpdateColumnsSize() const;
    void UpdateData() const;
    void Clear() const;

    void start() const;

public Q_SLOTS:
    void toggleSearchPanel() const;
    void onFilterChanged() const;
    void handleTableDoubleClicked( const QModelIndex &index ) const;
    void handleSessionsTableMenu(const QPoint &pos );

    void actionConsoleOpen() const;
    void actionExecuteCommand();
    void actionTasksBrowserOpen() const;
    void actionMarkActive() const;
    void actionMarkInactive() const;
    void actionItemColor() const;
    void actionTextColor() const;
    void actionColorReset() const;
    void actionAgentRemove();
    void actionConsoleDelete();
    void actionItemTag() const;
    void actionItemHide() const;
    void actionItemsShowAll() const;
};

#endif
