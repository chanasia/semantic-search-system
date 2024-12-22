package hugot

import "github.com/knights-analytics/hugot"

type HugotSessionManager struct {
	session *hugot.Session
}

func NewHugotSessionManager() *HugotSessionManager {
	return &HugotSessionManager{}
}

func (m *HugotSessionManager) Initialize() error {
	session, err := hugot.NewORTSession()
	if err != nil {
		return err
	}
	m.session = session
	return nil
}

func (m *HugotSessionManager) Destroy() error {
	if m.session != nil {
		return m.session.Destroy()
	}
	return nil
}

func (m *HugotSessionManager) GetSession() interface{} {
	return m.session
}
