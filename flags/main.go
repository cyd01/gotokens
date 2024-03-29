package flags

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func forgevar(cmd, name string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ToUpper(cmd+"_"+name), "-", "_"), ".", "_")
}
func formatname(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}

type ALIAS struct {
	shortName string
	longName  string
}

type Flags struct {
	flagName    string
	flagSet     *flag.FlagSet
	flagAliases []ALIAS
}

func NewFlag(name string) *Flags {
	f := &Flags{
		flagName: name,
		flagSet:  flag.NewFlagSet(name, flag.ContinueOnError),
	}
	f.SetUsage(f.defaultUsage)
	return f
}

func (f *Flags) AliasByShort(name string) string {
	for _, val := range f.flagAliases {
		if val.shortName == name {
			return val.longName
		}
	}
	return ""
}

func (f *Flags) AliasByLong(name string) string {
	for _, val := range f.flagAliases {
		if val.longName == name {
			return val.shortName
		}
	}
	return ""
}

func (f *Flags) Arg(i int) string {
	return f.flagSet.Arg(i)
}

func (f *Flags) Args() []string {
	return f.flagSet.Args()
}

func (f *Flags) SetUsage(fn func()) {
	f.flagSet.Usage = fn
}

func (f *Flags) Usage() {
	f.flagSet.Usage()
}

func (f *Flags) defaultUsage() {
	if f.flagName == "" {
		fmt.Fprintf(f.flagSet.Output(), "Usage:\n")
	} else {
		fmt.Fprintf(f.flagSet.Output(), "Usage of %s:\n", f.flagName)
	}
	f.PrintDefaults()
	os.Exit(0)
}

func (f *Flags) Parse(arguments []string) error {
	args := arguments
	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			if strings.HasPrefix(args[i], "-") {
				var v string
				if strings.HasPrefix(args[i], "--") {
					v = args[i][2:]
				} else {
					v = args[i][1:]
				}
				n := f.AliasByShort(v)
				if len(n) > 0 {
					args[i] = "-" + n
				}
			}
		}
	}
	return f.flagSet.Parse(args)
}

func (f *Flags) addAlias(longName, shortName string) {
	if len(shortName) > 1 {
		panic(f.flagName + " short name too long: " + shortName)
	}
	if len(shortName) == 1 {
		if len(f.AliasByShort(shortName)) > 0 {
			panic(f.flagName + " short flag redefined:" + shortName)
		}
		f.flagAliases = append(f.flagAliases, ALIAS{shortName: shortName, longName: longName})
	}
}

// -- string Value
type stringValue string

func newStringValue(val string, p *string) *stringValue {
	*p = val
	return (*stringValue)(p)
}

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() interface{} { return string(*s) }

func (s *stringValue) String() string { return string(*s) }

func isZeroValue(fl *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	typ := reflect.TypeOf(fl.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	return value == z.Interface().(flag.Value).String()
}

func (f *Flags) PrintDefaults() {
	f.flagSet.VisitAll(func(fl *flag.Flag) {
		var b strings.Builder
		fmt.Fprintf(&b, "  -%s", fl.Name) // Two spaces before -; see next two comments.
		if v := f.AliasByLong(fl.Name); len(v) > 0 {
			fmt.Fprintf(&b, ", -%s", v)
		}
		name, usage := flag.UnquoteUsage(fl)
		if len(name) > 0 {
			b.WriteString(" ")
			b.WriteString(name)
		}
		// Boolean flags of one ASCII letter are so common we
		// treat them specially, putting their usage on the same line.
		if b.Len() <= 4 { // space, space, '-', 'x'.
			b.WriteString("\t")
		} else {
			// Four spaces before the tab triggers good alignment
			// for both 4- and 8-space tab stops.
			b.WriteString("\n    \t")
		}
		b.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))

		if !isZeroValue(fl, fl.DefValue) {
			if _, ok := fl.Value.(*stringValue); ok {
				// put quotes on the value
				fmt.Fprintf(&b, " (default %q)", fl.DefValue)
			} else {
				fmt.Fprintf(&b, " (default %v)", fl.DefValue)
			}
		}
		fmt.Fprint(f.flagSet.Output(), b.String(), "\n")
	})
}

/*
 * Bool
 */
func (f *Flags) BoolP(longName, shortName string, value bool, usage string) *bool {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Bool(longName, value, usage)
}
func (f *Flags) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	f.BoolVar(p, name, value, usage)
	return p
}
func (f *Flags) BoolVarP(vr *bool, longName, shortName string, value bool, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.BoolVar(vr, longName, value, usage)
}
func (f *Flags) BoolVar(vr *bool, name string, value bool, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if strings.EqualFold(v, "1") || strings.EqualFold(v, "true") {
			val = true
		} else {
			val = false
		}
	}
	f.flagSet.BoolVar(vr, formatname(name), val, usage)
}

/*
 * Duration
 */
func (f *Flags) DurationP(longName, shortName string, value time.Duration, usage string) *time.Duration {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Duration(longName, value, usage)
}
func (f *Flags) Duration(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	f.DurationVar(p, name, value, usage)
	return p
}
func (f *Flags) DurationVarP(vr *time.Duration, longName, shortName string, value time.Duration, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.DurationVar(vr, longName, value, usage)
}
func (f *Flags) DurationVar(vr *time.Duration, name string, value time.Duration, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if vv, err := time.ParseDuration(v); err == nil {
			val = vv
		}
	}
	f.flagSet.DurationVar(vr, formatname(name), val, usage)
}

/*
 * Int
 */
func (f *Flags) IntP(longName, shortName string, value int, usage string) *int {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Int(longName, value, usage)
}
func (f *Flags) Int(name string, value int, usage string) *int {
	p := new(int)
	f.IntVar(p, name, value, usage)
	return p
}
func (f *Flags) IntVarP(vr *int, longName, shortName string, value int, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.IntVar(vr, longName, value, usage)
}
func (f *Flags) IntVar(vr *int, name string, value int, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if vv, err := strconv.Atoi(v); err == nil {
			val = vv
		}
	}
	f.flagSet.IntVar(vr, formatname(name), val, usage)
}

/*
 * Float64
 */
func (f *Flags) Float64P(longName, shortName string, value float64, usage string) *float64 {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Float64(longName, value, usage)
}
func (f *Flags) Float64(name string, value float64, usage string) *float64 {
	p := new(float64)
	f.Float64Var(p, name, value, usage)
	return p
}
func (f *Flags) Float64VarP(vr *float64, longName, shortName string, value float64, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.Float64Var(vr, longName, value, usage)
}
func (f *Flags) Float64Var(vr *float64, name string, value float64, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if vv, err := strconv.ParseFloat(v, 64); err == nil {
			val = vv
		}
	}
	f.flagSet.Float64Var(vr, formatname(name), val, usage)
}

/*
 * Int64
 */
func (f *Flags) Int64P(longName, shortName string, value int64, usage string) *int64 {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Int64(longName, value, usage)
}
func (f *Flags) Int64(name string, value int64, usage string) *int64 {
	p := new(int64)
	f.Int64Var(p, name, value, usage)
	return p
}
func (f *Flags) Int64VarP(vr *int64, longName, shortName string, value int64, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.Int64Var(vr, longName, value, usage)
}
func (f *Flags) Int64Var(vr *int64, name string, value int64, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if vv, err := strconv.ParseInt(v, 10, 64); err == nil {
			val = vv
		}
	}
	f.flagSet.Int64Var(vr, formatname(name), val, usage)
}

/*
 * String
 */
func (f *Flags) StringP(longName, shortName string, value string, usage string) *string {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.String(longName, value, usage)
}
func (f *Flags) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVar(p, name, value, usage)
	return p
}
func (f *Flags) StringVarP(vr *string, longName, shortName string, value string, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.StringVar(vr, longName, value, usage)
}
func (f *Flags) StringVar(vr *string, name string, value string, usage string) {
	val := value
	if v, err := f.getGlobalConfig(name); err == nil {
		val = v
	}
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		val = v
	}
	f.flagSet.StringVar(vr, formatname(name), val, usage)
}

func (f *Flags) NArg() int {
	return f.flagSet.NArg()
}

func (f *Flags) NFlag() int {
	return f.flagSet.NFlag()
}

func (f *Flags) Name() string {
	return f.flagSet.Name()
}

func (f *Flags) Parsed() bool {
	return f.flagSet.Parsed()
}

func (f *Flags) Set(name, value string) error {
	return f.flagSet.Set(name, value)
}

/*
 * Uint
 */
func (f *Flags) UintP(longName, shortName string, value uint, usage string) *uint {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Uint(longName, value, usage)
}
func (f *Flags) Uint(name string, value uint, usage string) *uint {
	p := new(uint)
	f.UintVar(p, name, value, usage)
	return p
}
func (f *Flags) UintVarP(vr *uint, longName, shortName string, value uint, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.UintVar(vr, longName, value, usage)
}
func (f *Flags) UintVar(vr *uint, name string, value uint, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if vv, err := strconv.ParseUint(v, 10, 64); err == nil {
			val = uint(vv)
		}
	}
	f.flagSet.UintVar(vr, formatname(name), val, usage)
}

/*
 * Uint64
 */
func (f *Flags) Uint64P(longName, shortName string, value uint64, usage string) *uint64 {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	return f.Uint64(longName, value, usage)
}
func (f *Flags) Uint64(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	f.Uint64Var(p, name, value, usage)
	return p
}
func (f *Flags) Uint64VarP(vr *uint64, longName, shortName string, value uint64, usage string) {
	if len(shortName) >= 1 {
		f.addAlias(formatname(longName), shortName)
	}
	f.Uint64Var(vr, longName, value, usage)
}
func (f *Flags) Uint64Var(vr *uint64, name string, value uint64, usage string) {
	val := value
	n := forgevar(f.flagName, name)
	if v, b := os.LookupEnv(n); b {
		if vv, err := strconv.ParseUint(v, 10, 64); err == nil {
			val = vv
		}
	}
	f.flagSet.Uint64Var(vr, formatname(name), val, usage)
}

// existsfile return if the given path exists and is a regular file
func existsfile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	if !fi.Mode().IsRegular() {
		return false
	}
	return true
}

func (f *Flags) getGlobalConfig(name string) (string, error) {
	val := ""
	filename := os.Getenv("HOME") + "/." + f.flagName + "/config.json"
	if !existsfile(filename) {
		return val, errors.New("Environment file not found")
	}
	var content []byte
	var err error = nil
	if content, err = ioutil.ReadFile(filename); err != nil {
		return val, errors.New("Can not read environment file")
	}
	data := make(map[string]interface{})
	if err = json.Unmarshal(content, &data); err != nil {
		return val, errors.New("Can not unmarshal environment file")
	}
	found := false
	for k, v := range data {
		if k == name {
			val = v.(string)
			found = true
		}
	}
	if found {
		return val, nil
	} else {
		return "", errors.New("Not found")
	}
}
