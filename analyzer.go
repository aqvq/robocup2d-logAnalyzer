package main

import (
	"bufio"
	"encoding/json"
	"github.com/spf13/cast"
	"io"
	"os"
	"strings"
)

type VectorStr struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type PlayerIDStr struct {
	Side string `json:"side"`
	Num  string `json:"num"`
}

type ArmStr struct {
	IsPointing bool      `json:"isPointing"`
	Dest       VectorStr `json:"destination"`
}

type FocusStr struct {
	IsFocusing bool        `json:"isFocusing"`
	Player     PlayerIDStr `json:"player"`
}

type BallTypeStr struct {
	Pos VectorStr `json:"position"`
	Vel VectorStr `json:"velocity"`
}

type ActionCountStr struct {
	Kick        string `json:"kick"`
	Dash        string `json:"dash"`
	Turn        string `json:"turn"`
	Catch       string `json:"catch"`
	Move        string `json:"move"`
	TurnNeck    string `json:"turnNeck"`
	ChangeView  string `json:"changeView"`
	Say         string `json:"say"`
	Tackle      string `json:"tackle"`
	Arm         string `json:"arm"`
	AttentionTo string `json:"attentionTo"`
}

type PlayerStr struct {
	ID              PlayerIDStr    `json:"id"`
	Type            string         `json:"typeId"`
	Flag            string         `json:"flag"`
	Pos             VectorStr      `json:"position"`
	Vel             VectorStr      `json:"velocity"`
	BodyAngle       VectorStr      `json:"bodyAngle"`
	Arm             ArmStr         `json:"arm,omitempty"`
	ViewMode        string         `json:"viewMode"`
	VisibleAngle    string         `json:"visibleAngle"`
	Stamina         string         `json:"stamina"`
	Effort          string         `json:"effort"`
	Recovery        string         `json:"recovery"`
	StaminaCapacity string         `json:"staminaCapacity"`
	Focus           FocusStr       `json:"focus,omitempty"`
	Counts          ActionCountStr `json:"counts"`
}

type CycleStr struct {
	ID       string        `json:"cycle"`
	Ball     BallTypeStr   `json:"ball"`
	Players  [22]PlayerStr `json:"players"`
	PlayMode string        `json:"playMode"`
	Score    [2]string     `json:"score"`
}

// AnalyzerStr 将整个文件转换成一行一行进行处理
func AnalyzerStr(source string, dest string, marshal bool, callback func(string)) {
	sourceFile, err := os.Open(source)
	if err != nil {
		panic("Error opening file")
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
		}
	}(sourceFile)
	sourceReader := bufio.NewReader(sourceFile)

	destFile, err := os.Create(dest)
	if err != nil {
		panic("Error creating file")
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {

		}
	}(destFile)
	destEncoder := json.NewEncoder(destFile)

	cycle := new(CycleStr)
	for {
		line, err := sourceReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			panic("Error reading string")
		}
		if err == io.EOF {
			break
		}
		arr := strings.Split(line, " ")
		processPlayModeStr(arr, cycle)
		processTeamStr(arr, cycle)
		if processShowStr(arr, cycle) {
			if marshal {
				data, err := json.MarshalIndent(cycle, "", "  ")
				if err != nil {
					panic("Error constructing json file")
				}
				_, err2 := destFile.Write(data)
				if err2 != nil {
					return
				}
			} else {
				err := destEncoder.Encode(cycle)
				if err != nil {
					panic("Error encoding")
				}
			}
		}
	}
	if callback != nil {
		callback(dest)
	}
}

// 处理每行的show数据
func processShowStr(arr []string, cycle *CycleStr) bool {
	//cycle := new(Cycle)
	offset := 0

	if arr[0] != "(show" {
		return false
	}

	cycle.ID = arr[1]
	cycle.Ball.Pos = VectorStr{arr[3], arr[4]}
	cycle.Ball.Vel = VectorStr{arr[5], TrimBracket(arr[6])}
	offset += 7

	for i := 0; i < 22; i++ {
		// (l或(r部分
		if strings.HasSuffix(arr[offset], "l") {
			cycle.Players[i].ID = PlayerIDStr{Side: "left", Num: TrimBracket(arr[offset+1])}
		} else {
			cycle.Players[i].ID = PlayerIDStr{Side: "right", Num: TrimBracket(arr[offset+1])}
		}
		offset += 2

		// 基本信息部分
		cycle.Players[i].Type = arr[offset]
		cycle.Players[i].Flag = arr[offset+1]
		cycle.Players[i].Pos = VectorStr{arr[offset+2], arr[offset+3]}
		cycle.Players[i].Vel = VectorStr{arr[offset+4], arr[offset+5]}
		cycle.Players[i].BodyAngle = VectorStr{arr[offset+6], arr[offset+7]}
		offset += 8

		// arm.pointing是可省略的
		if arr[offset] != "(v" {
			cycle.Players[i].Arm.IsPointing = true
			cycle.Players[i].Arm.Dest = VectorStr{arr[offset], arr[offset+1]}
			offset += 2
		} else {
			cycle.Players[i].Arm.IsPointing = false
			cycle.Players[i].Arm.Dest = VectorStr{"0", "0"}
		}

		// (v部分
		if arr[offset+1] == "h" {
			cycle.Players[i].ViewMode = "high"
		} else {
			cycle.Players[i].ViewMode = "low"
		}
		cycle.Players[i].VisibleAngle = TrimBracket(arr[offset+2])
		offset += 3

		// (s部分
		cycle.Players[i].Stamina = arr[offset+1]
		cycle.Players[i].Effort = arr[offset+2]
		cycle.Players[i].Recovery = arr[offset+3]
		cycle.Players[i].StaminaCapacity = TrimBracket(arr[offset+4])
		offset += 5

		// (f部分（可省略）
		if arr[offset] == "(f" {
			cycle.Players[i].Focus.IsFocusing = true
			if arr[offset+1] == "l" {
				cycle.Players[i].Focus.Player = PlayerIDStr{Side: "left", Num: TrimBracket(arr[offset+2])}
			} else {
				cycle.Players[i].Focus.Player = PlayerIDStr{Side: "right", Num: TrimBracket(arr[offset+2])}
			}
			offset += 3
		} else {
			cycle.Players[i].Focus.IsFocusing = false
			cycle.Players[i].Focus.Player = PlayerIDStr{Side: "", Num: "0"}
		}

		// (c部分
		cycle.Players[i].Counts = ActionCountStr{
			Kick:        arr[offset+1],
			Dash:        arr[offset+2],
			Turn:        arr[offset+3],
			Catch:       arr[offset+4],
			Move:        arr[offset+5],
			TurnNeck:    arr[offset+6],
			ChangeView:  arr[offset+7],
			Say:         arr[offset+8],
			Tackle:      arr[offset+9],
			Arm:         arr[offset+10],
			AttentionTo: TrimBracket(arr[offset+11])}
		offset += 12
	}
	return true
}

// 处理(team数据
func processTeamStr(arr []string, cycle *CycleStr) {
	if arr[0] != "(team" {
		return
	}
	offset := 4
	for {
		if offset-4 >= len(cycle.Score) {
			break
		}
		cycle.Score[offset-4] = TrimBracket(arr[offset])
		offset += 1
	}
}

// 处理(playmode数据
func processPlayModeStr(arr []string, cycle *CycleStr) {
	if arr[0] != "(playmode" {
		return
	}
	cycle.PlayMode = TrimBracket(arr[2])
}

type Vector struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type PlayerID struct {
	Side string `json:"side"`
	Num  uint8  `json:"num"`
}

type Arm struct {
	IsPointing bool   `json:"isPointing"`
	Dest       Vector `json:"destination"`
}

type Focus struct {
	IsFocusing bool     `json:"isFocusing"`
	Player     PlayerID `json:"player"`
}

type BallType struct {
	Pos Vector `json:"position"`
	Vel Vector `json:"velocity"`
}

type ActionCount struct {
	Kick        uint `json:"kick"`
	Dash        uint `json:"dash"`
	Turn        uint `json:"turn"`
	Catch       uint `json:"catch"`
	Move        uint `json:"move"`
	TurnNeck    uint `json:"turnNeck"`
	ChangeView  uint `json:"changeView"`
	Say         uint `json:"say"`
	Tackle      uint `json:"tackle"`
	Arm         uint `json:"arm"`
	AttentionTo uint `json:"attentionTo"`
}

type Player struct {
	ID              PlayerID    `json:"id"`
	Type            uint8       `json:"typeId"`
	Flag            uint32      `json:"flag"`
	Pos             Vector      `json:"position"`
	Vel             Vector      `json:"velocity"`
	BodyAngle       Vector      `json:"bodyAngle"`
	Arm             Arm         `json:"arm,omitempty"`
	ViewMode        string      `json:"viewMode"`
	VisibleAngle    uint8       `json:"visibleAngle"`
	Stamina         float32     `json:"stamina"`
	Effort          float32     `json:"effort"`
	Recovery        float32     `json:"recovery"`
	StaminaCapacity float32     `json:"staminaCapacity"`
	Focus           Focus       `json:"focus,omitempty"`
	Counts          ActionCount `json:"counts"`
}

type Cycle struct {
	ID       uint16     `json:"cycle"`
	Ball     BallType   `json:"ball"`
	Players  [22]Player `json:"players"`
	PlayMode string     `json:"playMode"`
	Score    [2]uint16  `json:"score"`
}

// Analyzer 将整个文件转换成一行一行进行处理
func Analyzer(source string, dest string, marshal bool, callback func(string)) {
	sourceFile, err := os.Open(source)
	if err != nil {
		panic("Error opening file")
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
		}
	}(sourceFile)
	sourceReader := bufio.NewReader(sourceFile)

	destFile, err := os.Create(dest)
	if err != nil {
		panic("Error creating file")
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {

		}
	}(destFile)
	destEncoder := json.NewEncoder(destFile)

	cycle := new(Cycle)
	for {
		line, err := sourceReader.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil && err != io.EOF {
			panic("Error reading string")
		}
		if err == io.EOF {
			break
		}
		arr := strings.Split(line, " ")
		processPlayMode(arr, cycle)
		processTeam(arr, cycle)
		if processShow(arr, cycle) {
			if marshal {
				data, err := json.MarshalIndent(cycle, "", "  ")
				if err != nil {
					panic("Error constructing json file")
				}
				_, err2 := destFile.Write(data)
				if err2 != nil {
					return
				}
			} else {
				err := destEncoder.Encode(cycle)
				if err != nil {
					panic("Error encoding")
				}
			}
		}
	}
	if callback != nil {
		callback(dest)
	}
}

// 处理每行的show数据
func processShow(arr []string, cycle *Cycle) bool {
	//cycle := new(Cycle)
	if arr[0] != "(show" {
		return false
	}

	offset := 0
	cycle.ID = cast.ToUint16(arr[1])
	cycle.Ball.Pos = Vector{cast.ToFloat32(arr[3]), cast.ToFloat32(arr[4])}
	cycle.Ball.Vel = Vector{cast.ToFloat32(arr[5]), cast.ToFloat32(TrimBracket(arr[6]))}
	cycle.Players = [22]Player{}
	offset += 7

	for i := 0; i < 22; i++ {
		// (l或(r部分
		if strings.HasSuffix(arr[offset], "l") {
			cycle.Players[i].ID = PlayerID{Side: "left", Num: cast.ToUint8(TrimBracket(arr[offset+1]))}
		} else {
			cycle.Players[i].ID = PlayerID{Side: "right", Num: cast.ToUint8(TrimBracket(arr[offset+1]))}
		}
		offset += 2

		// 基本信息部分
		cycle.Players[i].Type = cast.ToUint8(arr[offset])
		cycle.Players[i].Flag = cast.ToUint32(arr[offset+1])
		cycle.Players[i].Pos = Vector{cast.ToFloat32(arr[offset+2]), cast.ToFloat32(arr[offset+3])}
		cycle.Players[i].Vel = Vector{cast.ToFloat32(arr[offset+4]), cast.ToFloat32(arr[offset+5])}
		cycle.Players[i].BodyAngle = Vector{cast.ToFloat32(arr[offset+6]), cast.ToFloat32(arr[offset+7])}
		offset += 8

		// arm.pointing是可省略的
		if arr[offset] != "(v" {
			cycle.Players[i].Arm.IsPointing = true
			cycle.Players[i].Arm.Dest = Vector{cast.ToFloat32(arr[offset]), cast.ToFloat32(arr[offset+1])}
			offset += 2
		}

		// (v部分
		if arr[offset+1] == "h" {
			cycle.Players[i].ViewMode = "high"
		} else {
			cycle.Players[i].ViewMode = "low"
		}
		cycle.Players[i].VisibleAngle = cast.ToUint8(TrimBracket(arr[offset+2]))
		offset += 3

		// (s部分
		cycle.Players[i].Stamina = cast.ToFloat32(arr[offset+1])
		cycle.Players[i].Effort = cast.ToFloat32(arr[offset+2])
		cycle.Players[i].Recovery = cast.ToFloat32(arr[offset+3])
		cycle.Players[i].StaminaCapacity = cast.ToFloat32(TrimBracket(arr[offset+4]))
		offset += 5

		// (f部分（可省略）
		if arr[offset] == "(f" {
			cycle.Players[i].Focus.IsFocusing = true
			if arr[offset+1] == "l" {
				cycle.Players[i].Focus.Player = PlayerID{Side: "left", Num: cast.ToUint8(TrimBracket(arr[offset+2]))}
			} else {
				cycle.Players[i].Focus.Player = PlayerID{Side: "right", Num: cast.ToUint8(TrimBracket(arr[offset+2]))}
			}
			offset += 3
		} else {
			cycle.Players[i].Focus.IsFocusing = false
		}

		// (c部分
		cycle.Players[i].Counts = ActionCount{
			Kick:        cast.ToUint(arr[offset+1]),
			Dash:        cast.ToUint(arr[offset+2]),
			Turn:        cast.ToUint(arr[offset+3]),
			Catch:       cast.ToUint(arr[offset+4]),
			Move:        cast.ToUint(arr[offset+5]),
			TurnNeck:    cast.ToUint(arr[offset+6]),
			ChangeView:  cast.ToUint(arr[offset+7]),
			Say:         cast.ToUint(arr[offset+8]),
			Tackle:      cast.ToUint(arr[offset+9]),
			Arm:         cast.ToUint(arr[offset+10]),
			AttentionTo: cast.ToUint(TrimBracket(arr[offset+11]))}
		offset += 12
	}
	return true
}

// 处理(team数据
func processTeam(arr []string, cycle *Cycle) {
	if arr[0] != "(team" {
		return
	}
	offset := 4
	for {
		if offset-4 >= len(arr) {
			break
		}
		if !strings.HasSuffix(arr[offset], ")") {
			cycle.Score[offset-4] = cast.ToUint16(arr[offset])
			offset += 1
		} else {
			cycle.Score[offset-4] = cast.ToUint16(TrimBracket(arr[offset]))
			break
		}
	}
}

// 处理(playmode数据
func processPlayMode(arr []string, cycle *Cycle) {
	if arr[0] != "(playmode" {
		return
	}
	cycle.PlayMode = TrimBracket(arr[2])
}
