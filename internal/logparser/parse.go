package logparser

import (
    "bufio"
    "strings"
    "kla/internal/model"
)

// ParseLogs reads logs from a string and returns a list of LogEntry instances
func ParseLogs(logData string) ([]model.LogEntry, error) {
    var logs []model.LogEntry
    scanner := bufio.NewScanner(strings.NewReader(logData))
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(strings.ToLower(line), "err") || strings.Contains(strings.ToLower(line), "error") {
            logs = append(logs, model.LogEntry{
                Message: line,
            })
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return logs, nil
}