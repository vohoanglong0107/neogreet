package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogoDim(t *testing.T) {
	testCases := []struct {
		name           string
		logo           string
		expectedHeight int
		expectedWidth  int
	}{
		{"fedora logo",
			`[blue]             .:dddl:.
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
    .:odl:.`, 12, 23},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			width, heigth := getDim(tc.logo)
			require.Equal(t, tc.expectedWidth, width)
			require.Equal(t, tc.expectedHeight, heigth)
		})
	}
}
