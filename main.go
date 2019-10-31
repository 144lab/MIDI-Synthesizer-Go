package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gordonklaus/portaudio"
	"github.com/rakyll/portmidi"
)

const (
	// SampleRate ...
	SampleRate = 44100
	// MaxPolyphony ...
	MaxPolyphony = 16
)

// Params ...
type Params struct {
	attack       float32
	decay        float32
	sustainLevel float32
	sustainRate  float32
	sustain      float32
	release      float32
	form         []float32
	formGain     float32
	formRate     float32
}

// Note ...
type Note struct {
	num   uint8
	vel   uint8
	on    bool
	top   bool
	gain  float32
	phase float32
}

var (
	waveForm0 = []float32{0.0, 0.049067674327418015, 0.0980171403295606, 0.14673047445536175, 0.19509032201612825, 0.24298017990326387, 0.29028467725446233, 0.33688985339222005, 0.3826834323650898, 0.4275550934302821, 0.47139673682599764, 0.5141027441932217, 0.5555702330196022, 0.5956993044924334, 0.6343932841636455, 0.6715589548470183, 0.7071067811865475, 0.740951125354959, 0.7730104533627369, 0.8032075314806448, 0.8314696123025451, 0.8577286100002721, 0.8819212643483549, 0.9039892931234433, 0.9238795325112867, 0.9415440651830208, 0.9569403357322089, 0.970031253194544, 0.9807852804032304, 0.989176509964781, 0.9951847266721968, 0.9987954562051724, 1.0, 0.9987954562051724, 0.9951847266721969, 0.989176509964781, 0.9807852804032304, 0.970031253194544, 0.9569403357322089, 0.9415440651830208, 0.9238795325112867, 0.9039892931234434, 0.881921264348355, 0.8577286100002721, 0.8314696123025455, 0.8032075314806449, 0.7730104533627371, 0.740951125354959, 0.7071067811865476, 0.6715589548470186, 0.6343932841636455, 0.5956993044924335, 0.5555702330196022, 0.5141027441932218, 0.4713967368259978, 0.42755509343028203, 0.38268343236508984, 0.3368898533922203, 0.29028467725446233, 0.24298017990326404, 0.19509032201612858, 0.1467304744553618, 0.09801714032956084, 0.04906767432741797, 1.2246467991473532e-16, -0.049067674327417724, -0.09801714032956059, -0.14673047445536158, -0.19509032201612836, -0.24298017990326382, -0.29028467725446216, -0.3368898533922201, -0.38268343236508967, -0.4275550934302818, -0.47139673682599764, -0.5141027441932216, -0.555570233019602, -0.5956993044924332, -0.6343932841636453, -0.6715589548470184, -0.7071067811865475, -0.7409511253549588, -0.7730104533627367, -0.803207531480645, -0.8314696123025452, -0.857728610000272, -0.8819212643483549, -0.9039892931234431, -0.9238795325112865, -0.9415440651830208, -0.9569403357322088, -0.970031253194544, -0.9807852804032303, -0.9891765099647809, -0.9951847266721969, -0.9987954562051724, -1.0, -0.9987954562051724, -0.9951847266721969, -0.9891765099647809, -0.9807852804032304, -0.970031253194544, -0.9569403357322089, -0.9415440651830209, -0.9238795325112866, -0.9039892931234433, -0.881921264348355, -0.8577286100002722, -0.8314696123025456, -0.8032075314806453, -0.7730104533627369, -0.7409511253549592, -0.7071067811865477, -0.6715589548470187, -0.6343932841636459, -0.5956993044924332, -0.5555702330196022, -0.5141027441932219, -0.4713967368259979, -0.42755509343028253, -0.3826834323650904, -0.33688985339222, -0.29028467725446244, -0.24298017990326418, -0.19509032201612872, -0.1467304744553624, -0.09801714032956052, -0.0490676743274180}
	waveForm1 = []float32{
		0.02228, 0.02765, 0.02053, -0.01519, -0.06787, -0.06072, -0.04724,
		-0.05765, -0.01910, 0.03289, 0.08273, 0.13940, 0.11706, 0.06689,
		0.03338, -0.05285, -0.13366, -0.16491, -0.20506, -0.20248, -0.16497,
		-0.13637, -0.04612, 0.05436, 0.10167, 0.11837, 0.09184, 0.03813,
		-0.03075, -0.12342, -0.19155, -0.24263, -0.33331, -0.35835, -0.29186,
		-0.18098, -0.04185, 0.05620, 0.10544, 0.15068, 0.16912, 0.14566,
		0.10187, 0.02552, -0.05956, -0.14947, -0.19638, -0.19449, -0.15328,
		-0.11134, -0.10183, -0.06497, -0.00833, 0.04250, 0.09087, 0.12332,
		0.14033, 0.14816, 0.14266, 0.14905, 0.18481, 0.20491, 0.19229,
		0.17313, 0.14034, 0.12052, 0.12403, 0.12324, 0.10882, 0.11918,
		0.12034, 0.13370, 0.17856, 0.22018, 0.25724, 0.26085, 0.24341,
		0.22923, 0.21964, 0.18427, 0.15720, 0.13564, 0.08770, 0.05586,
		0.06215, 0.11416, 0.18126, 0.17765, 0.13364, 0.10718, 0.09003,
		0.07549, 0.04007, 0.02016, 0.03321, 0.02301, 0.04958, 0.09636,
		0.13982, 0.21460, 0.21760, 0.14192, 0.01850, -0.13288, -0.29390,
		-0.50247, -0.67149, -0.70905, -0.70960, -0.62262, -0.41466, -0.19885,
		0.04470, 0.26506, 0.36032, 0.39179, 0.38315, 0.24268, 0.06933,
		-0.06901, -0.24582, -0.34780, -0.35623, -0.37080, -0.29164, -0.17080,
		-0.11854, -0.04797,
	}
	toneMap = []float32{
		16.351597831287414, 17.323914436054505,
		18.354047994837977, 19.445436482630058,
		20.601722307054366, 21.826764464562746,
		23.12465141947715, 24.499714748859326,
		25.956543598746574, 27.5,
		29.13523509488062, 30.86770632850775,
		32.70319566257483, 34.64782887210901,
		36.70809598967594, 38.890872965260115,
		41.20344461410875, 43.653528929125486,
		46.2493028389543, 48.999429497718666,
		51.91308719749314, 55.0,
		58.27047018976124, 61.7354126570155,
		65.40639132514966, 69.29565774421802,
		73.41619197935188, 77.78174593052023,
		82.4068892282175, 87.30705785825097,
		92.4986056779086, 97.99885899543733,
		103.82617439498628, 110.0,
		116.54094037952248, 123.47082531403103,
		130.8127826502993, 138.59131548843604,
		146.8323839587038, 155.56349186104046,
		164.81377845643496, 174.61411571650194,
		184.9972113558172, 195.99771799087463,
		207.65234878997256, 220.0,
		233.08188075904496, 246.94165062806206,
		261.6255653005986, 277.1826309768721,
		293.6647679174076, 311.1269837220809,
		329.6275569128699, 349.2282314330039,
		369.9944227116344, 391.99543598174927,
		415.3046975799451, 440.0,
		466.1637615180899, 493.8833012561241,
		523.2511306011972, 554.3652619537442,
		587.3295358348151, 622.2539674441618,
		659.2551138257398, 698.4564628660078,
		739.9888454232688, 783.9908719634985,
		830.6093951598903, 880.0,
		932.3275230361799, 987.7666025122483,
		1046.5022612023945, 1108.7305239074883,
		1174.6590716696303, 1244.5079348883237,
		1318.5102276514797, 1396.9129257320155,
		1479.9776908465376, 1567.981743926997,
		1661.2187903197805, 1760.0,
		1864.6550460723597, 1975.533205024496,
		2093.004522404789, 2217.4610478149766,
		2349.31814333926, 2489.0158697766474,
		2637.02045530296, 2793.825851464031,
		2959.955381693075, 3135.9634878539946,
		3322.437580639561, 3520.0,
		3729.3100921447194, 3951.066410048992,
		4186.009044809578, 4434.922095629953,
		4698.63628667852, 4978.031739553295,
		5274.04091060592, 5587.651702928062,
		5919.91076338615, 6271.926975707989,
		6644.875161279122, 7040.0,
		7458.620184289437, 7902.132820097988,
		8372.018089619156, 8869.844191259906,
		9397.272573357044, 9956.06347910659,
		10548.081821211836, 11175.303405856126,
		11839.8215267723, 12543.853951415975,
		13289.750322558246, 14080.0,
		14917.240368578874, 15804.265640195976,
		16744.036179238312, 17739.688382519813,
		18794.54514671409, 19912.12695821318,
		21096.16364242367, 22350.606811712252,
		23679.6430535446, 25087.70790283195,
		26579.50064511649, 28160.0,
	}

	paramsList = []*Params{
		&Params{
			attack:       10. / SampleRate,
			decay:        5. / SampleRate,
			sustainLevel: 0.9,
			sustainRate:  9.0,
			sustain:      0.44 / SampleRate,
			release:      2. / SampleRate,
			form:         waveForm0,
			formGain:     0.5,
			formRate:     1.0,
		},
		&Params{
			attack:       10. / SampleRate,
			decay:        5. / SampleRate,
			sustainLevel: 0.7,
			sustainRate:  9.0,
			sustain:      0.4 / SampleRate,
			release:      2. / SampleRate,
			form:         waveForm1,
			formGain:     0.8,
			formRate:     0.3,
		},
	}
	paramsIndex = 0
	notes       = [MaxPolyphony]Note{}
	control     = [128]uint8{}
	pitch       = float32(0)
	mutex       sync.RWMutex
)

func noteOn(num, vel uint8) {
	minGain := float32(1.0)
	var note Note
	index := 0
	phaseReset := true
	mutex.Lock()
	defer mutex.Unlock()
	for i, n := range notes {
		if n.num == num {
			note = n
			index = i
			phaseReset = false
			break
		}
		if minGain > n.gain {
			minGain = n.gain
			note = n
			index = i
		}
	}
	note.num = num
	note.vel = vel
	note.on = true
	if phaseReset {
		note.phase = 0.0
	}
	notes[index] = note
}

func noteOff(num, vel uint8) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, n := range notes {
		if n.num == num {
			n.on = false
			notes[i] = n
			return
		}
	}
}

func env(n *Note) float32 {
	mutex.Lock()
	defer mutex.Unlock()
	if n.gain == 0.0 && n.vel == 0 {
		n.top = false
		n.phase = 0.0
		return 0.0
	}
	topLevel := float32(n.vel) / 127.0
	params := paramsList[paramsIndex]
	if n.on {
		if !n.top {
			n.gain += params.attack
			if n.gain > topLevel {
				n.top = true
				n.gain = topLevel
			}
		} else {
			if n.gain > params.sustainLevel*topLevel {
				n.gain -= params.decay
			} else {
				n.gain -= params.sustain
			}
			if n.gain < 0.0 {
				n.gain = 0.0
				n.vel = 0.0
				n.phase = 0
				n.num = 0
			}
		}
	} else {
		n.top = false
		if n.gain > params.sustainLevel*topLevel {
			n.gain -= params.decay
		} else {
			release := params.release / (1.0 + float32(control[64])*params.sustainRate/127.0)
			n.gain -= release
		}
		if n.gain < 0.0 {
			n.gain = 0.0
			n.vel = 0.0
			n.phase = 0
			n.num = 0
		}
	}
	return n.gain
}

func operate(n *Note, f float32) float32 {
	mutex.RLock()
	defer mutex.RUnlock()
	if n.gain == 0.0 && n.vel == 0 {
		return 0.0
	}
	params := paramsList[paramsIndex]
	formLen := len(params.form)
	formLenF := float32(formLen)
	i := int(formLenF * n.phase)
	p := n.phase*formLenF - float32(i)
	v := params.form[i%formLen] * (1 - p)
	v += params.form[(i+1)%formLen] * (p)
	v *= n.gain * params.formGain
	n.phase += f / SampleRate * params.formRate
	n.phase = n.phase - float32(int(n.phase))
	if v > 1.0 {
		return 1.0
	}
	if v < -1.0 {
		return -1.0
	}
	return v
}

func procSample(out [][]float32) {
	for j := range out[0] {
		var v float32
		for i := range notes {
			env(&notes[i])
			// ピッチベンド処理（最大全音＝半音２つ分シフト）
			n := toneMap[notes[i].num]
			m := toneMap[notes[i].num+2]
			p := pitch
			if pitch < 0 {
				p = -1 * pitch
				i2 := int(notes[i].num) - 2
				if i2 < 0 {
					i2 = 0
				}
				m = toneMap[i2]
			}
			f := n*(1-p) + m*(p)
			v += 0.8 * operate(&notes[i], f)
		}
		if v > 1.0 {
			v = 1.0
		}
		if v < -1.0 {
			v = -1.0
		}
		out[0][j] = v
		out[1][j] = v
	}
}

func main() {
	portmidi.Initialize()
	defer portmidi.Terminate()
	log.Println("devices:", portmidi.CountDevices())
	deviceID := portmidi.DefaultInputDeviceID()
	in, err := portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	ch := in.Listen()
	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 2, SampleRate, 0, procSample)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
	if err := stream.Start(); err != nil {
		log.Fatal(err)
	}
	defer stream.Stop()
	for {
		select {
		case msg := <-ch:
			switch msg.Status {
			case 0x90: // note on
				fmt.Printf("\rNoteON  message bytes = %02x %02x\x1b[K", msg.Data1, msg.Data2)
				noteOn(uint8(msg.Data1), uint8(msg.Data2))
			case 0x80: //note off
				fmt.Printf("\rNoteOFF message bytes = %02x %02x\x1b[K", msg.Data1, msg.Data2)
				noteOff(uint8(msg.Data1), uint8(msg.Data2))
			case 224: // pitch
				// ピッチベンドを-1.0〜1.0の範囲に起こす
				p := (float32(msg.Data2) - 64.0) / 63.0
				if p < -1.0 {
					p = -1.0
				}
				fmt.Printf("\rPitch Bend message = %+5.2f\x1b[K", p)
				mutex.Lock()
				pitch = p
				mutex.Unlock()
			case 176: // control change
				mutex.Lock()
				last := control[uint8(msg.Data1)&127]
				control[uint8(msg.Data1)&127] = uint8(msg.Data2) & 127
				switch msg.Data1 {
				case 1:
					if last < 127 && msg.Data2 == 127 {
						paramsIndex++
						if paramsIndex >= len(paramsList) {
							paramsIndex = len(paramsList) - 1
						}
						fmt.Printf("\rChange WaveForm: %02x\x1b[K\n", paramsIndex)
					}
				case 2:
					if last < 127 && msg.Data2 == 127 {
						paramsIndex--
						if paramsIndex < 0 {
							paramsIndex = 0
						}
						fmt.Printf("\rChange WaveForm: %02x\x1b[K\n", paramsIndex)
					}
				default:
					fmt.Printf("\rControl Change message bytes = %02x %02x\x1b[K", msg.Data1, msg.Data2)
				}
				mutex.Unlock()
			default:
				fmt.Printf("%d message bytes = %02x %02x\x1b[K\n", msg.Status, msg.Data1, msg.Data2)
			}
		}
	}
}
