
Datatypes
=========


.. go:function:: func NewDateDuration(months int32, days int32) DateDuration

    NewDateDuration returns a new DateDuration


.. go:type:: type DateDuration struct {\
        // contains filtered or unexported fields\
    }

    DateDuration represents the elapsed time between two dates in a fuzzy human
    way.


.. go:method:: func (dd DateDuration) MarshalText() ([]byte, error)

    MarshalText returns dd marshaled as text.


.. go:method:: func (dd DateDuration) String() string

    


.. go:method:: func (dd *DateDuration) UnmarshalText(b []byte) error

    UnmarshalText unmarshals bytes into \*dd.


.. go:function:: func DurationFromNanoseconds(d time.Duration) Duration

    DurationFromNanoseconds creates a Duration represented as microseconds
    from a `time.Duration <https://pkg.go.dev/time>`_ represented as nanoseconds.


.. go:type:: type Duration int64

    Duration represents the elapsed time between two instants
    as an int64 microsecond count.


.. go:method:: func (d Duration) AsNanoseconds() (time.Duration, error)

    AsNanoseconds returns `time.Duration <https://pkg.go.dev/time>`_ represented as nanoseconds,
    after transforming from Duration microsecond representation.
    Returns an error if the Duration is too long and would cause an overflow of
    the internal int64 representation.


.. go:method:: func (d Duration) String() string

    


.. go:function:: func NewLocalDate(year int, month time.Month, day int) LocalDate

    NewLocalDate returns a new LocalDate


.. go:type:: type LocalDate struct {\
        // contains filtered or unexported fields\
    }

    LocalDate is a date without a time zone.
    `docs/stdlib/datetime#type::cal::local_date <https://www.edgedb.com/docs/stdlib/datetime#type::cal::local_date>`_


.. go:method:: func (d LocalDate) MarshalText() ([]byte, error)

    MarshalText returns d marshaled as text.


.. go:method:: func (d LocalDate) String() string

    


.. go:method:: func (d *LocalDate) UnmarshalText(b []byte) error

    UnmarshalText unmarshals bytes into \*d.


.. go:function:: func NewLocalDateTime(\
        year int, month time.Month, day, hour, minute, second, microsecond int,\
    ) LocalDateTime

    NewLocalDateTime returns a new LocalDateTime


.. go:type:: type LocalDateTime struct {\
        // contains filtered or unexported fields\
    }

    LocalDateTime is a date and time without timezone.
    `docs/stdlib/datetime#type::cal::local_datetime <https://www.edgedb.com/docs/stdlib/datetime#type::cal::local_datetime>`_


.. go:method:: func (dt LocalDateTime) MarshalText() ([]byte, error)

    MarshalText returns dt marshaled as text.


.. go:method:: func (dt LocalDateTime) String() string

    


.. go:method:: func (dt *LocalDateTime) UnmarshalText(b []byte) error

    UnmarshalText unmarshals bytes into \*dt.


.. go:function:: func NewLocalTime(hour, minute, second, microsecond int) LocalTime

    NewLocalTime returns a new LocalTime


.. go:type:: type LocalTime struct {\
        // contains filtered or unexported fields\
    }

    LocalTime is a time without a time zone.
    `docs/stdlib/datetime#type::cal::local_time <https://www.edgedb.com/docs/stdlib/datetime#type::cal::local_time>`_


.. go:method:: func (t LocalTime) MarshalText() ([]byte, error)

    MarshalText returns t marshaled as text.


.. go:method:: func (t LocalTime) String() string

    


.. go:method:: func (t *LocalTime) UnmarshalText(b []byte) error

    UnmarshalText unmarshals bytes into \*t.


.. go:type:: type Memory int64

    Memory represents memory in bytes.


.. go:method:: func (m Memory) MarshalText() ([]byte, error)

    MarshalText returns m marshaled as text.


.. go:method:: func (m Memory) String() string

    


.. go:method:: func (m *Memory) UnmarshalText(b []byte) error

    UnmarshalText unmarshals bytes into \*m.


.. go:type:: type Optional struct {\
        // contains filtered or unexported fields\
    }

    Optional represents a shape field that is not required.
    Optional is embedded in structs to make them optional. For example:
    
    .. code-block:: go
    
        type User struct {
            edgedb.Optional
            Name string `edgedb:"name"`
        }


.. go:method:: func (o *Optional) Missing() bool

    Missing returns true if the value is missing.


.. go:method:: func (o *Optional) SetMissing(missing bool)

    SetMissing sets the structs missing status. true means missing and false
    means present.


.. go:method:: func (o *Optional) Unset()

    Unset marks the value as missing


.. go:function:: func NewOptionalBigInt(v *big.Int) OptionalBigInt

    NewOptionalBigInt is a convenience function for creating an OptionalBigInt
    with its value set to v.


.. go:type:: type OptionalBigInt struct {\
        // contains filtered or unexported fields\
    }

    OptionalBigInt is an optional \*big.Int. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalBigInt) Get() (*big.Int, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalBigInt) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalBigInt) Set(val *big.Int)

    Set sets the value.


.. go:method:: func (o *OptionalBigInt) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalBigInt) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalBool(v bool) OptionalBool

    NewOptionalBool is a convenience function for creating an OptionalBool with
    its value set to v.


.. go:type:: type OptionalBool struct {\
        // contains filtered or unexported fields\
    }

    OptionalBool is an optional bool. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalBool) Get() (bool, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalBool) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalBool) Set(val bool)

    Set sets the value.


.. go:method:: func (o *OptionalBool) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalBool) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalBytes(v []byte) OptionalBytes

    NewOptionalBytes is a convenience function for creating an OptionalBytes
    with its value set to v.


.. go:type:: type OptionalBytes struct {\
        // contains filtered or unexported fields\
    }

    OptionalBytes is an optional []byte. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalBytes) Get() ([]byte, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalBytes) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalBytes) Set(val []byte)

    Set sets the value.


.. go:method:: func (o *OptionalBytes) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalBytes) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalDateDuration(v DateDuration) OptionalDateDuration

    NewOptionalDateDuration is a convenience function for creating an
    OptionalDateDuration with its value set to v.


.. go:type:: type OptionalDateDuration struct {\
        // contains filtered or unexported fields\
    }

    OptionalDateDuration is an optional DateDuration. Optional types
    must be used for out parameters when a shape field is not required.


.. go:method:: func (o *OptionalDateDuration) Get() (DateDuration, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalDateDuration) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalDateDuration) Set(val DateDuration)

    Set sets the value.


.. go:method:: func (o *OptionalDateDuration) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalDateDuration) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalDateTime(v time.Time) OptionalDateTime

    NewOptionalDateTime is a convenience function for creating an
    OptionalDateTime with its value set to v.


.. go:type:: type OptionalDateTime struct {\
        // contains filtered or unexported fields\
    }

    OptionalDateTime is an optional time.Time.  Optional types must be used for
    out parameters when a shape field is not required.


.. go:method:: func (o OptionalDateTime) Get() (time.Time, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalDateTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalDateTime) Set(val time.Time)

    Set sets the value.


.. go:method:: func (o *OptionalDateTime) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalDateTime) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalDuration(v Duration) OptionalDuration

    NewOptionalDuration is a convenience function for creating an
    OptionalDuration with its value set to v.


.. go:type:: type OptionalDuration struct {\
        // contains filtered or unexported fields\
    }

    OptionalDuration is an optional Duration. Optional types must be used for
    out parameters when a shape field is not required.


.. go:method:: func (o OptionalDuration) Get() (Duration, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalDuration) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalDuration) Set(val Duration)

    Set sets the value.


.. go:method:: func (o *OptionalDuration) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalDuration) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalFloat32(v float32) OptionalFloat32

    NewOptionalFloat32 is a convenience function for creating an OptionalFloat32
    with its value set to v.


.. go:type:: type OptionalFloat32 struct {\
        // contains filtered or unexported fields\
    }

    OptionalFloat32 is an optional float32. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalFloat32) Get() (float32, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalFloat32) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalFloat32) Set(val float32)

    Set sets the value.


.. go:method:: func (o *OptionalFloat32) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalFloat32) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalFloat64(v float64) OptionalFloat64

    NewOptionalFloat64 is a convenience function for creating an OptionalFloat64
    with its value set to v.


.. go:type:: type OptionalFloat64 struct {\
        // contains filtered or unexported fields\
    }

    OptionalFloat64 is an optional float64. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalFloat64) Get() (float64, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalFloat64) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalFloat64) Set(val float64)

    Set sets the value.


.. go:method:: func (o *OptionalFloat64) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalFloat64) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalInt16(v int16) OptionalInt16

    NewOptionalInt16 is a convenience function for creating an OptionalInt16
    with its value set to v.


.. go:type:: type OptionalInt16 struct {\
        // contains filtered or unexported fields\
    }

    OptionalInt16 is an optional int16. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalInt16) Get() (int16, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalInt16) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalInt16) Set(val int16)

    Set sets the value.


.. go:method:: func (o *OptionalInt16) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalInt16) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalInt32(v int32) OptionalInt32

    NewOptionalInt32 is a convenience function for creating an OptionalInt32
    with its value set to v.


.. go:type:: type OptionalInt32 struct {\
        // contains filtered or unexported fields\
    }

    OptionalInt32 is an optional int32. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalInt32) Get() (int32, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalInt32) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalInt32) Set(val int32)

    Set sets the value.


.. go:method:: func (o *OptionalInt32) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalInt32) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalInt64(v int64) OptionalInt64

    NewOptionalInt64 is a convenience function for creating an OptionalInt64
    with its value set to v.


.. go:type:: type OptionalInt64 struct {\
        // contains filtered or unexported fields\
    }

    OptionalInt64 is an optional int64. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalInt64) Get() (int64, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalInt64) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalInt64) Set(val int64) *OptionalInt64

    Set sets the value.


.. go:method:: func (o *OptionalInt64) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalInt64) Unset() *OptionalInt64

    Unset marks the value as missing.


.. go:function:: func NewOptionalLocalDate(v LocalDate) OptionalLocalDate

    NewOptionalLocalDate is a convenience function for creating an
    OptionalLocalDate with its value set to v.


.. go:type:: type OptionalLocalDate struct {\
        // contains filtered or unexported fields\
    }

    OptionalLocalDate is an optional LocalDate. Optional types must be used for
    out parameters when a shape field is not required.


.. go:method:: func (o OptionalLocalDate) Get() (LocalDate, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalLocalDate) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalLocalDate) Set(val LocalDate)

    Set sets the value.


.. go:method:: func (o *OptionalLocalDate) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalLocalDate) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalLocalDateTime(v LocalDateTime) OptionalLocalDateTime

    NewOptionalLocalDateTime is a convenience function for creating an
    OptionalLocalDateTime with its value set to v.


.. go:type:: type OptionalLocalDateTime struct {\
        // contains filtered or unexported fields\
    }

    OptionalLocalDateTime is an optional LocalDateTime. Optional types must be
    used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalLocalDateTime) Get() (LocalDateTime, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalLocalDateTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalLocalDateTime) Set(val LocalDateTime)

    Set sets the value.


.. go:method:: func (o *OptionalLocalDateTime) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalLocalDateTime) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalLocalTime(v LocalTime) OptionalLocalTime

    NewOptionalLocalTime is a convenience function for creating an
    OptionalLocalTime with its value set to v.


.. go:type:: type OptionalLocalTime struct {\
        // contains filtered or unexported fields\
    }

    OptionalLocalTime is an optional LocalTime. Optional types must be used for
    out parameters when a shape field is not required.


.. go:method:: func (o OptionalLocalTime) Get() (LocalTime, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalLocalTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalLocalTime) Set(val LocalTime)

    Set sets the value.


.. go:method:: func (o *OptionalLocalTime) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalLocalTime) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalMemory(v Memory) OptionalMemory

    NewOptionalMemory is a convenience function for creating an
    OptionalMemory with its value set to v.


.. go:type:: type OptionalMemory struct {\
        // contains filtered or unexported fields\
    }

    OptionalMemory is an optional Memory. Optional types must be used for
    out parameters when a shape field is not required.


.. go:method:: func (o OptionalMemory) Get() (Memory, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalMemory) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalMemory) Set(val Memory)

    Set sets the value.


.. go:method:: func (o *OptionalMemory) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalMemory) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeDateTime(v RangeDateTime) OptionalRangeDateTime

    NewOptionalRangeDateTime is a convenience function for creating an
    OptionalRangeDateTime with its value set to v.


.. go:type:: type OptionalRangeDateTime struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeDateTime is an optional RangeDateTime. Optional
    types must be used for out parameters when a shape field is not required.


.. go:method:: func (o *OptionalRangeDateTime) Get() (RangeDateTime, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o *OptionalRangeDateTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeDateTime) Set(val RangeDateTime)

    Set sets the value.


.. go:method:: func (o *OptionalRangeDateTime) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeDateTime) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeFloat32(v RangeFloat32) OptionalRangeFloat32

    NewOptionalRangeFloat32 is a convenience function for creating an
    OptionalRangeFloat32 with its value set to v.


.. go:type:: type OptionalRangeFloat32 struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeFloat32 is an optional RangeFloat32. Optional
    types must be used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRangeFloat32) Get() (RangeFloat32, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRangeFloat32) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeFloat32) Set(val RangeFloat32)

    Set sets the value.


.. go:method:: func (o *OptionalRangeFloat32) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeFloat32) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeFloat64(v RangeFloat64) OptionalRangeFloat64

    NewOptionalRangeFloat64 is a convenience function for creating an
    OptionalRangeFloat64 with its value set to v.


.. go:type:: type OptionalRangeFloat64 struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeFloat64 is an optional RangeFloat64. Optional
    types must be used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRangeFloat64) Get() (RangeFloat64, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRangeFloat64) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeFloat64) Set(val RangeFloat64)

    Set sets the value.


.. go:method:: func (o *OptionalRangeFloat64) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeFloat64) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeInt32(v RangeInt32) OptionalRangeInt32

    NewOptionalRangeInt32 is a convenience function for creating an
    OptionalRangeInt32 with its value set to v.


.. go:type:: type OptionalRangeInt32 struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeInt32 is an optional RangeInt32. Optional types must be used
    for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRangeInt32) Get() (RangeInt32, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRangeInt32) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeInt32) Set(val RangeInt32)

    Set sets the value.


.. go:method:: func (o *OptionalRangeInt32) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeInt32) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeInt64(v RangeInt64) OptionalRangeInt64

    NewOptionalRangeInt64 is a convenience function for creating an
    OptionalRangeInt64 with its value set to v.


.. go:type:: type OptionalRangeInt64 struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeInt64 is an optional RangeInt64. Optional
    types must be used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRangeInt64) Get() (RangeInt64, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRangeInt64) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeInt64) Set(val RangeInt64)

    Set sets the value.


.. go:method:: func (o *OptionalRangeInt64) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeInt64) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeLocalDate(v RangeLocalDate) OptionalRangeLocalDate

    NewOptionalRangeLocalDate is a convenience function for creating an
    OptionalRangeLocalDate with its value set to v.


.. go:type:: type OptionalRangeLocalDate struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeLocalDate is an optional RangeLocalDate. Optional types must be
    used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRangeLocalDate) Get() (RangeLocalDate, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRangeLocalDate) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeLocalDate) Set(val RangeLocalDate)

    Set sets the value.


.. go:method:: func (o *OptionalRangeLocalDate) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeLocalDate) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRangeLocalDateTime(\
        v RangeLocalDateTime,\
    ) OptionalRangeLocalDateTime

    NewOptionalRangeLocalDateTime is a convenience function for creating an
    OptionalRangeLocalDateTime with its value set to v.


.. go:type:: type OptionalRangeLocalDateTime struct {\
        // contains filtered or unexported fields\
    }

    OptionalRangeLocalDateTime is an optional RangeLocalDateTime. Optional
    types must be used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRangeLocalDateTime) Get() (RangeLocalDateTime, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRangeLocalDateTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRangeLocalDateTime) Set(val RangeLocalDateTime)

    Set sets the value.


.. go:method:: func (o *OptionalRangeLocalDateTime) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRangeLocalDateTime) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalRelativeDuration(v RelativeDuration) OptionalRelativeDuration

    NewOptionalRelativeDuration is a convenience function for creating an
    OptionalRelativeDuration with its value set to v.


.. go:type:: type OptionalRelativeDuration struct {\
        // contains filtered or unexported fields\
    }

    OptionalRelativeDuration is an optional RelativeDuration. Optional types
    must be used for out parameters when a shape field is not required.


.. go:method:: func (o OptionalRelativeDuration) Get() (RelativeDuration, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalRelativeDuration) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalRelativeDuration) Set(val RelativeDuration)

    Set sets the value.


.. go:method:: func (o *OptionalRelativeDuration) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalRelativeDuration) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalStr(v string) OptionalStr

    NewOptionalStr is a convenience function for creating an OptionalStr with
    its value set to v.


.. go:type:: type OptionalStr struct {\
        // contains filtered or unexported fields\
    }

    OptionalStr is an optional string. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalStr) Get() (string, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalStr) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalStr) Set(val string)

    Set sets the value.


.. go:method:: func (o *OptionalStr) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o.


.. go:method:: func (o *OptionalStr) Unset()

    Unset marks the value as missing.


.. go:function:: func NewOptionalUUID(v UUID) OptionalUUID

    NewOptionalUUID is a convenience function for creating an OptionalUUID with
    its value set to v.


.. go:type:: type OptionalUUID struct {\
        // contains filtered or unexported fields\
    }

    OptionalUUID is an optional UUID. Optional types must be used for out
    parameters when a shape field is not required.


.. go:method:: func (o OptionalUUID) Get() (UUID, bool)

    Get returns the value and a boolean indicating if the value is present.


.. go:method:: func (o OptionalUUID) MarshalJSON() ([]byte, error)

    MarshalJSON returns o marshaled as json.


.. go:method:: func (o *OptionalUUID) Set(val UUID)

    Set sets the value.


.. go:method:: func (o *OptionalUUID) UnmarshalJSON(bytes []byte) error

    UnmarshalJSON unmarshals bytes into \*o


.. go:method:: func (o *OptionalUUID) Unset()

    Unset marks the value as missing.


.. go:function:: func NewRangeDateTime(\
        lower, upper OptionalDateTime,\
        incLower, incUpper bool,\
    ) RangeDateTime

    NewRangeDateTime creates a new RangeDateTime value.


.. go:type:: type RangeDateTime struct {\
        // contains filtered or unexported fields\
    }

    RangeDateTime is an interval of time.Time values.


.. go:method:: func (r RangeDateTime) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeDateTime) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeDateTime) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeDateTime) Lower() OptionalDateTime

    Lower returns the lower bound.


.. go:method:: func (r RangeDateTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeDateTime) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeDateTime) Upper() OptionalDateTime

    Upper returns the upper bound.


.. go:function:: func NewRangeFloat32(\
        lower, upper OptionalFloat32,\
        incLower, incUpper bool,\
    ) RangeFloat32

    NewRangeFloat32 creates a new RangeFloat32 value.


.. go:type:: type RangeFloat32 struct {\
        // contains filtered or unexported fields\
    }

    RangeFloat32 is an interval of float32 values.


.. go:method:: func (r RangeFloat32) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeFloat32) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeFloat32) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeFloat32) Lower() OptionalFloat32

    Lower returns the lower bound.


.. go:method:: func (r RangeFloat32) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeFloat32) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeFloat32) Upper() OptionalFloat32

    Upper returns the upper bound.


.. go:function:: func NewRangeFloat64(\
        lower, upper OptionalFloat64,\
        incLower, incUpper bool,\
    ) RangeFloat64

    NewRangeFloat64 creates a new RangeFloat64 value.


.. go:type:: type RangeFloat64 struct {\
        // contains filtered or unexported fields\
    }

    RangeFloat64 is an interval of float64 values.


.. go:method:: func (r RangeFloat64) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeFloat64) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeFloat64) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeFloat64) Lower() OptionalFloat64

    Lower returns the lower bound.


.. go:method:: func (r RangeFloat64) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeFloat64) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeFloat64) Upper() OptionalFloat64

    Upper returns the upper bound.


.. go:function:: func NewRangeInt32(\
        lower, upper OptionalInt32,\
        incLower, incUpper bool,\
    ) RangeInt32

    NewRangeInt32 creates a new RangeInt32 value.


.. go:type:: type RangeInt32 struct {\
        // contains filtered or unexported fields\
    }

    RangeInt32 is an interval of int32 values.


.. go:method:: func (r RangeInt32) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeInt32) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeInt32) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeInt32) Lower() OptionalInt32

    Lower returns the lower bound.


.. go:method:: func (r RangeInt32) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeInt32) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeInt32) Upper() OptionalInt32

    Upper returns the upper bound.


.. go:function:: func NewRangeInt64(\
        lower, upper OptionalInt64,\
        incLower, incUpper bool,\
    ) RangeInt64

    NewRangeInt64 creates a new RangeInt64 value.


.. go:type:: type RangeInt64 struct {\
        // contains filtered or unexported fields\
    }

    RangeInt64 is an interval of int64 values.


.. go:method:: func (r RangeInt64) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeInt64) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeInt64) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeInt64) Lower() OptionalInt64

    Lower returns the lower bound.


.. go:method:: func (r RangeInt64) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeInt64) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeInt64) Upper() OptionalInt64

    Upper returns the upper bound.


.. go:function:: func NewRangeLocalDate(\
        lower, upper OptionalLocalDate,\
        incLower, incUpper bool,\
    ) RangeLocalDate

    NewRangeLocalDate creates a new RangeLocalDate value.


.. go:type:: type RangeLocalDate struct {\
        // contains filtered or unexported fields\
    }

    RangeLocalDate is an interval of LocalDate values.


.. go:method:: func (r RangeLocalDate) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeLocalDate) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeLocalDate) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeLocalDate) Lower() OptionalLocalDate

    Lower returns the lower bound.


.. go:method:: func (r RangeLocalDate) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeLocalDate) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeLocalDate) Upper() OptionalLocalDate

    Upper returns the upper bound.


.. go:function:: func NewRangeLocalDateTime(\
        lower, upper OptionalLocalDateTime,\
        incLower, incUpper bool,\
    ) RangeLocalDateTime

    NewRangeLocalDateTime creates a new RangeLocalDateTime value.


.. go:type:: type RangeLocalDateTime struct {\
        // contains filtered or unexported fields\
    }

    RangeLocalDateTime is an interval of LocalDateTime values.


.. go:method:: func (r RangeLocalDateTime) Empty() bool

    Empty returns true if the range is empty.


.. go:method:: func (r RangeLocalDateTime) IncLower() bool

    IncLower returns true if the lower bound is inclusive.


.. go:method:: func (r RangeLocalDateTime) IncUpper() bool

    IncUpper returns true if the upper bound is inclusive.


.. go:method:: func (r RangeLocalDateTime) Lower() OptionalLocalDateTime

    Lower returns the lower bound.


.. go:method:: func (r RangeLocalDateTime) MarshalJSON() ([]byte, error)

    MarshalJSON returns r marshaled as json.


.. go:method:: func (r *RangeLocalDateTime) UnmarshalJSON(data []byte) error

    UnmarshalJSON unmarshals bytes into \*r.


.. go:method:: func (r RangeLocalDateTime) Upper() OptionalLocalDateTime

    Upper returns the upper bound.


.. go:function:: func NewRelativeDuration(\
        months, days int32,\
        microseconds int64,\
    ) RelativeDuration

    NewRelativeDuration returns a new RelativeDuration


.. go:type:: type RelativeDuration struct {\
        // contains filtered or unexported fields\
    }

    RelativeDuration represents the elapsed time between two instants in a fuzzy
    human way.


.. go:method:: func (rd RelativeDuration) MarshalText() ([]byte, error)

    MarshalText returns rd marshaled as text.


.. go:method:: func (rd RelativeDuration) String() string

    


.. go:method:: func (rd *RelativeDuration) UnmarshalText(b []byte) error

    UnmarshalText unmarshals bytes into \*rd.


.. go:function:: func ParseUUID(s string) (UUID, error)

    ParseUUID parses s into a UUID or returns an error.


.. go:type:: type UUID [16]byte

    UUID is a universally unique identifier
    `docs/stdlib/uuid <https://www.edgedb.com/docs/stdlib/uuid>`_


.. go:method:: func (id UUID) MarshalText() ([]byte, error)

    MarshalText returns the id as a byte string.


.. go:method:: func (id UUID) String() string

    


.. go:method:: func (id *UUID) UnmarshalText(b []byte) error

    UnmarshalText unmarshals the id from a string.