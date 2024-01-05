package main

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Setup(t testing.TB) (string, func()) {
	tmpfile, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name(), func() {
		os.Remove(tmpfile.Name())
	}
}

func TestLogo(t *testing.T) {
	configPath, tearDown := Setup(t)
	t.Cleanup(tearDown)

	logo := `[blue]             .:dddl:.
            OWMKOOXMWd
           KMMc    xMMc
           MMM.     WW:
           MMM.
    oxOOOo MMM0OOk.
  0MMKxdd: MMMkddc.
 XM0'      MMM.
 MMo       MMW.
 0MNc.   .xMMd
  dNMWXXXWM0:
    .:odl:.
`
	content := `
  logo:
    template: |
      [blue]             .:dddl:.
                  OWMKOOXMWd
                 KMMc    xMMc
                 MMM.     WW:
                 MMM.
          oxOOOo MMM0OOk.
        0MMKxdd: MMMkddc.
       XM0'      MMM.
       MMo       MMW.
       0MNc.   .xMMd
        dNMWXXXWM0:
          .:odl:.
  `
	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	config := NewConfig(configPath)
	require.Equal(t, logo, config.getLogo())

}

type MockInfo struct {
	mock.Mock
}

func (mockInfo *MockInfo) CPU() string {
	args := mockInfo.Called()
	return args.String(0)
}

func (mockInfo *MockInfo) GPU() string {
	args := mockInfo.Called()
	return args.String(0)
}

func (mockInfo *MockInfo) Memory() string {
	args := mockInfo.Called()
	return args.String(0)
}

func (mockInfo *MockInfo) Disk() string {
	args := mockInfo.Called()
	return args.String(0)
}

func (mockInfo *MockInfo) OS() string {
	args := mockInfo.Called()
	return args.String(0)
}

func TestInfo(t *testing.T) {
	configPath, tearDown := Setup(t)
	t.Cleanup(tearDown)

	info := `cpu: Core I5
gpu: Nvidia
`
	content := `
  info:
    template: |
      cpu: {{ .CPU }}
      gpu: {{ .GPU }}
  `
	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
	config := NewConfig(configPath)
	mockInfo := new(MockInfo)
	mockInfo.On("CPU").Return("Core I5")
	mockInfo.On("GPU").Return("Nvidia")
	require.Equal(t, info, config.getInfo(mockInfo))

}
