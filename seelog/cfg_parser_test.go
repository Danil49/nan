// Copyright (c) 2012 - Cloud Instruments Co., Ltd.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package seelog

import (
	"fmt"
	//"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var re *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func getTestFileName(testName, postfix string) string {
	if len(postfix) != 0 {
		return strings.ToLower(re.ReplaceAllString(testName, "_")) + "_" + postfix + "_test.log"
	}
	return strings.ToLower(re.ReplaceAllString(testName, "_")) + "_test.log"
}

var parserTests []parserTest

type parserTest struct {
	testName      string
	config        string
	expected      *logConfig //interface{}
	errorExpected bool
}

func getParserTests() []parserTest {
	if parserTests == nil {
		parserTests = make([]parserTest, 0)

		testName := "Simple file output"
		testLogFileName := getTestFileName(testName, "")
		testConfig := `
		<seelog>
			<outputs>
				<file path="` + testLogFileName + `"/>
			</outputs>
		</seelog>
		`
		testExpected := new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testfileWriter, _ := newFileWriter(testLogFileName)
		testHeadSplitter, _ := newSplitDispatcher(defaultformatter, []interface{}{testfileWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Filter dispatcher"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<filter levels="debug, info, critical">
					<file path="` + testLogFileName + `"/>
				</filter>
			</outputs>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testfileWriter, _ = newFileWriter(testLogFileName)
		testFilter, _ := newFilterDispatcher(defaultformatter, []interface{}{testfileWriter}, DebugLvl, InfoLvl, CriticalLvl)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testFilter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Console writer"
		testConfig = `
		<seelog type="sync">
			<outputs>
				<console />
			</outputs>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ := newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Smtp writer"
		testConfig = `
<seelog>
	<outputs>
		<smtp senderaddress="sa" sendername="sn"  hostname="hn" hostport="123" username="un" password="up">
			<recipient address="ra1"/>
			<recipient address="ra2"/>
			<recipient address="ra3"/>
			<cacertdirpath path="cacdp1"/>
			<cacertdirpath path="cacdp2"/>
		</smtp>
	</outputs>
</seelog>
		`

		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testSmtpWriter := newSmtpWriter(
			"sa",
			"sn",
			[]string{"ra1", "ra2", "ra3"},
			"hn",
			"123",
			"un",
			"up",
			[]string{"cacdp1", "cacdp2"},
		)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testSmtpWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Default output"
		testConfig = `
		<seelog type="sync"/>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Asyncloop behavior"
		testConfig = `
		<seelog type="asyncloop"/>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Asynctimer behavior"
		testConfig = `
		<seelog type="asynctimer" asyncinterval="101"/>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncTimerloggerTypeFromString
		testExpected.LoggerData = asyncTimerLoggerData{101}
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Rolling file writer size"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<rollingfile type="size" filename="` + testLogFileName + `" maxsize="100" maxrolls="5" />
			</outputs>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ := newRollingFileWriterSize(testLogFileName, rollingArchiveNone, "", 100, 5)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testrollingFileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Rolling file writer archive zip"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<rollingfile type="size" filename="` + testLogFileName + `" maxsize="100" maxrolls="5" archivetype="zip"/>
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ = newRollingFileWriterSize(testLogFileName, rollingArchiveZip, "log.zip", 100, 5)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testrollingFileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Rolling file writer archive zip with specified path"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<rollingfile type="size" filename="` + testLogFileName + `" maxsize="100" maxrolls="5" archivetype="zip" archivepath="test.zip"/>
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ = newRollingFileWriterSize(testLogFileName, rollingArchiveZip, "test.zip", 100, 5)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testrollingFileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Rolling file writer archive none"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<rollingfile type="size" filename="` + testLogFileName + `" maxsize="100" maxrolls="5" archivetype="none"/>
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ = newRollingFileWriterSize(testLogFileName, rollingArchiveNone, "", 100, 5)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testrollingFileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Rolling file writer date"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<rollingfile type="date" filename="` + testLogFileName + `" datepattern="2006-01-02T15:04:05Z07:00" />
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ = newRollingFileWriterDate(testLogFileName, rollingArchiveNone, "", "2006-01-02T15:04:05Z07:00")
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testrollingFileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Buffered writer"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<buffered size="100500" flushperiod="100">
					<rollingfile type="date" filename="` + testLogFileName + `" datepattern="2006-01-02T15:04:05Z07:00" />
				</buffered>
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ = newRollingFileWriterDate(testLogFileName, rollingArchiveNone, "", "2006-01-02T15:04:05Z07:00")
		testbufferedWriter, _ := newBufferedWriter(testrollingFileWriter, 100500, 100)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testbufferedWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Inner splitter output"
		testLogFileName1 := getTestFileName(testName, "1")
		testLogFileName2 := getTestFileName(testName, "2")
		testLogFileName3 := getTestFileName(testName, "3")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<file path="` + testLogFileName1 + `"/>
				<splitter>
					<file path="` + testLogFileName2 + `"/>
					<file path="` + testLogFileName3 + `"/>
				</splitter>
			</outputs>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testfileWriter1, _ := newFileWriter(testLogFileName2)
		testfileWriter2, _ := newFileWriter(testLogFileName3)
		testInnerSplitter, _ := newSplitDispatcher(defaultformatter, []interface{}{testfileWriter1, testfileWriter2})
		testfileWriter, _ = newFileWriter(testLogFileName1)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testfileWriter, testInnerSplitter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Format"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs formatid="dateFormat">
				<file path="` + testLogFileName + `"/>
			</outputs>
			<formats>
				<format id="dateFormat" format="%Level %Msg %File" />
			</formats>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testfileWriter, _ = newFileWriter(testLogFileName)
		testFormat, _ := newFormatter("%Level %Msg %File")
		testHeadSplitter, _ = newSplitDispatcher(testFormat, []interface{}{testfileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Format2"
		testLogFileName = getTestFileName(testName, "")
		testLogFileName1 = getTestFileName(testName, "1")
		testConfig = `
		<seelog type="sync">
			<outputs formatid="format1">
				<file path="` + testLogFileName + `"/>
				<file formatid="format2" path="` + testLogFileName1 + `"/>
			</outputs>
			<formats>
				<format id="format1" format="%Level %Msg %File" />
				<format id="format2" format="%l %Msg" />
			</formats>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testfileWriter, _ = newFileWriter(testLogFileName)
		testfileWriter1, _ = newFileWriter(testLogFileName1)
		testFormat1, _ := newFormatter("%Level %Msg %File")
		testFormat2, _ := newFormatter("%l %Msg")
		formattedWriter, _ := newFormattedWriter(testfileWriter1, testFormat2)
		testHeadSplitter, _ = newSplitDispatcher(testFormat1, []interface{}{testfileWriter, formattedWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Minlevel = warn"
		testConfig = `<seelog minlevel="warn"/>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(WarnLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Maxlevel = trace"
		testConfig = `<seelog maxlevel="trace"/>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, TraceLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Level between info and error"
		testConfig = `<seelog minlevel="info" maxlevel="error"/>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(InfoLvl, ErrorLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Off with minlevel"
		testConfig = `<seelog minlevel="off"/>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newOffConstraints()
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Off with levels"
		testConfig = `<seelog levels="off"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Levels list"
		testConfig = `<seelog levels="debug, info, critical"/>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newListConstraints([]LogLevel{
			DebugLvl, InfoLvl, CriticalLvl})
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = asyncLooploggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Errors #1"
		testConfig = `<seelog minlevel="debug" minlevel="trace"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #2"
		testConfig = `<seelog minlevel="error" maxlevel="debug"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #3"
		testConfig = `<seelog maxlevel="debug" maxlevel="trace"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #4"
		testConfig = `<seelog maxlevel="off"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #5"
		testConfig = `<seelog minlevel="off" maxlevel="trace"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #6"
		testConfig = `<seelog minlevel="warn" maxlevel="error" levels="debug"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #7"
		testConfig = `<not_seelog/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #8"
		testConfig = `<seelog levels="warn, debug, test"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #9"
		testConfig = `<seelog levels=""/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #10"
		testConfig = `<seelog levels="off" something="abc"/>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #11"
		testConfig = `<seelog><output/></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #12"
		testConfig = `<seelog><outputs/><outputs/></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #13"
		testConfig = `<seelog><exceptions/></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #14"
		testConfig = `<seelog><formats/></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #15"
		testConfig = `<seelog><outputs><splitter/></outputs></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #16"
		testConfig = `<seelog><outputs><filter/></outputs></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #17"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `<seelog><outputs><file path="` + testLogFileName + `"><something/></file></outputs></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #18"
		testConfig = `<seelog><outputs><buffered size="100500" flushperiod="100"/></outputs></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #19"
		testConfig = `<seelog><outputs></outputs></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Exceptions: restricting"
		testConfig =
			`
		<seelog type="sync">
			<exceptions>
				<exception funcpattern="Test*" filepattern="someFile.go" minlevel="off"/>
			</exceptions>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		listConstraint, _ := newOffConstraints()
		exception, _ := newLogLevelException("Test*", "someFile.go", listConstraint)
		testExpected.Exceptions = []*logLevelException{exception}
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Exceptions: allowing #1"
		testConfig =
			`
		<seelog type="sync" levels="error">
			<exceptions>
				<exception filepattern="testfile.go" minlevel="trace"/>
			</exceptions>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newListConstraints([]LogLevel{ErrorLvl})
		minMaxConstraint, _ := newMinMaxConstraints(TraceLvl, CriticalLvl)
		exception, _ = newLogLevelException("*", "testfile.go", minMaxConstraint)
		testExpected.Exceptions = []*logLevelException{exception}
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Exceptions: allowing #2"
		testConfig = `
		<seelog type="sync" levels="off">
			<exceptions>
				<exception filepattern="testfile.go" minlevel="warn"/>
			</exceptions>
		</seelog>
		`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newOffConstraints()
		minMaxConstraint, _ = newMinMaxConstraints(WarnLvl, CriticalLvl)
		exception, _ = newLogLevelException("*", "testfile.go", minMaxConstraint)
		testExpected.Exceptions = []*logLevelException{exception}
		testconsoleWriter, _ = newConsoleWriter()
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testconsoleWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Predefined formats"
		formatId := predefinedPrefix + "xml-debug-short"
		testConfig = `
		<seelog type="sync">
			<outputs formatid="` + formatId + `">
				<console />
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testconsoleWriter, _ = newConsoleWriter()
		testFormat, _ = predefinedFormats[formatId]
		testHeadSplitter, _ = newSplitDispatcher(testFormat, []interface{}{testconsoleWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Predefined formats redefine"
		testLogFileName = getTestFileName(testName, "")
		formatId = predefinedPrefix + "xml-debug-short"
		testConfig = `
		<seelog type="sync">
			<outputs formatid="` + formatId + `">
				<file path="` + testLogFileName + `"/>
			</outputs>
			<formats>
				<format id="` + formatId + `" format="%Level %Msg %File" />
			</formats>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testfileWriter, _ = newFileWriter(testLogFileName)
		testFormat, _ = newFormatter("%Level %Msg %File")
		testHeadSplitter, _ = newSplitDispatcher(testFormat, []interface{}{testfileWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Conn writer 1"
		testConfig = `
		<seelog type="sync">
			<outputs>
				<conn net="tcp" addr=":8888" />
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testConnWriter := newConnWriter("tcp", ":8888", false)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testConnWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Conn writer 2"
		testConfig = `
		<seelog type="sync">
			<outputs>
				<conn net="tcp" addr=":8888" reconnectonmsg="true" />
			</outputs>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testConnWriter = newConnWriter("tcp", ":8888", true)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{testConnWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

		testName = "Errors #11"
		testConfig = `
		<seelog type="sync"><exceptions>
				<exception filepattern="testfile.go" minlevel="trace"/>
				<exception filepattern="testfile.go" minlevel="warn"/>
		</exceptions></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #12"
		testConfig = `
		<seelog type="sync"><exceptions>
				<exception filepattern="!@+$)!!%&@(^$" minlevel="trace"/>
		</exceptions></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #13"
		testConfig = `
		<seelog type="sync"><exceptions>
				<exception filepattern="*" minlevel="unknown"/>
		</exceptions></seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #14"
		testConfig = `
		<seelog type="sync" levels=???off???>
			<exceptions>
				<exception filepattern="testfile.go" minlevel="off"/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #15"
		testConfig = `
		<seelog type="sync" levels=???trace???>
			<exceptions>
				<exception filepattern="testfile.go" levels="trace"/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #16"
		testConfig = `
		<seelog type="sync" minlevel=???trace???>
			<exceptions>
				<exception filepattern="testfile.go" minlevel="trace"/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #17"
		testConfig = `
		<seelog type="sync" minlevel=???trace???>
			<exceptions>
				<exception filepattern="testfile.go" minlevel="warn"/>
			</exceptions>
			<exceptions>
				<exception filepattern="testfile.go" minlevel="warn"/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #18"
		testConfig = `
		<seelog type="sync" minlevel=???trace???>
			<exceptions>
				<exception filepattern="testfile.go"/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #19"
		testConfig = `
		<seelog type="sync" minlevel=???trace???>
			<exceptions>
				<exception minlevel="warn"/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #20"
		testConfig = `
		<seelog type="sync" minlevel=???trace???>
			<exceptions>
				<exception/>
			</exceptions>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #21"
		testConfig = `
		<seelog>
			<outputs>
				<splitter>
				</splitter>
			</outputs>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #22"
		testConfig = `
		<seelog type="sync">
			<outputs>
				<filter levels="debug, info, critical">

				</filter>
			</outputs>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #23"
		testConfig = `
		<seelog type="sync">
			<outputs>
				<buffered size="100500" flushperiod="100">

				</buffered>
			</outputs>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #24"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<buffered size="100500" flushperiod="100">
					<rollingfile type="date" filename="` + testLogFileName + `" datepattern="2006-01-02T15:04:05Z07:00" formatid="testFormat"/>
				</buffered>
			</outputs>
			<formats>
				<format id="testFormat" format="%Level %Msg %File 123" />
			</formats>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #25"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<outputs>
					<file path="` + testLogFileName + `"/>
				</outputs>
				<outputs>
					<file path="` + testLogFileName + `"/>
				</outputs>
			</outputs>
		</seelog>
		`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Errors #26"
		testConfig = `
		<seelog type="sync">
			<outputs>
				<conn net="tcp" addr=":8888" reconnectonmsg="true1" />
			</outputs>
		</seelog>`
		parserTests = append(parserTests, parserTest{testName, testConfig, nil, true})

		testName = "Buffered writer same formatid override"
		testLogFileName = getTestFileName(testName, "")
		testConfig = `
		<seelog type="sync">
			<outputs>
				<buffered size="100500" flushperiod="100" formatid="testFormat">
					<rollingfile type="date" filename="` + testLogFileName + `" datepattern="2006-01-02T15:04:05Z07:00" formatid="testFormat"/>
				</buffered>
			</outputs>
			<formats>
				<format id="testFormat" format="%Level %Msg %File 123" />
			</formats>
		</seelog>`
		testExpected = new(logConfig)
		testExpected.Constraints, _ = newMinMaxConstraints(TraceLvl, CriticalLvl)
		testExpected.Exceptions = nil
		testrollingFileWriter, _ = newRollingFileWriterDate(testLogFileName, rollingArchiveNone, "", "2006-01-02T15:04:05Z07:00")
		testbufferedWriter, _ = newBufferedWriter(testrollingFileWriter, 100500, 100)
		testFormat, _ = newFormatter("%Level %Msg %File 123")
		formattedWriter, _ = newFormattedWriter(testbufferedWriter, testFormat)
		testHeadSplitter, _ = newSplitDispatcher(defaultformatter, []interface{}{formattedWriter})
		testExpected.LogType = syncloggerTypeFromString
		testExpected.RootDispatcher = testHeadSplitter
		parserTests = append(parserTests, parserTest{testName, testConfig, testExpected, false})

	}

	return parserTests
}

// Temporary solution: compare by string identity.
func configsAreEqual(conf1 *logConfig, conf2 interface{}) bool {
	if conf1 == nil {
		return conf2 == nil
	}
	if conf2 == nil {
		return conf1 == nil
	}

	// logConfig, ok := conf2 //.(*logConfig)
	// if !ok {
	// 	return false
	// }

	return fmt.Sprintf("%s", conf1) == fmt.Sprintf("%s", conf2) //logConfig)
}

func testLogFileFilter(fn string) bool {
	return ".log" == filepath.Ext(fn)
}

func cleanupAfterCfgTest(t *testing.T) {
	toDel, err := getDirFilePaths(".", testLogFileFilter, true)
	if nil != err {
		t.Fatal("Cannot list files in test directory!")
	}

	for _, p := range toDel {
		err = tryRemoveFile(p)
		if nil != err {
			t.Errorf("Cannot remove file %s in test directory: %s", p, err.Error())
		}
	}
}

func parseTest(test parserTest, t *testing.T) {
	conf, err := configFromReader(strings.NewReader(test.config))
	if /*err != nil &&*/ conf != nil && conf.RootDispatcher != nil {
		defer func() {
			if err = conf.RootDispatcher.Close(); err != nil {
				t.Errorf("\n----ERROR while closing root dispatcher in %s test: %s", test.testName, err)
			}
		}()
	}

	if (err != nil) != test.errorExpected {
		t.Errorf("\n----ERROR in %s:\nConfig: %s\n* Expected error:%t. Got error: %t\n",
			test.testName, test.config, test.errorExpected, (err != nil))
		if err != nil {
			t.Logf("%s\n", err.Error())
		}
		return
	}

	if err == nil && !configsAreEqual(conf, test.expected) {
		t.Errorf("\n----ERROR in %s:\nConfig: %s\n* Expected: %s. \n* Got: %s\n",
			test.testName, test.config, test.expected, conf)
	}
}

func TestParser(t *testing.T) {
	defer cleanupAfterCfgTest(t)

	for _, test := range getParserTests() {
		parseTest(test, t)
	}
}
