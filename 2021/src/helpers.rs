use std::{error, fmt, io, num};

#[derive(Debug)]
pub enum ParseError {
    ParseIntError(num::ParseIntError),
    ParseStringError(String),
    IOError(io::Error),
}

impl fmt::Display for ParseError {
    fn fmt(self: &Self, f: &mut fmt::Formatter) -> Result<(), fmt::Error> {
        write!(f, "{:?}", self)
    }
}

impl error::Error for ParseError {}

impl From<num::ParseIntError> for ParseError {
    fn from(e: num::ParseIntError) -> Self {
        ParseError::ParseIntError(e)
    }
}

impl From<String> for ParseError {
    fn from(e: String) -> Self {
        ParseError::ParseStringError(e)
    }
}

impl From<&str> for ParseError {
    fn from(e: &str) -> Self {
        ParseError::ParseStringError(e.into())
    }
}

impl From<io::Error> for ParseError {
    fn from(e: io::Error) -> Self {
        ParseError::IOError(e)
    }
}

pub fn split_fields<const N: usize>(s: &str) -> Result<[&str; N], ParseError> {
    s.split_whitespace()
        .collect::<Vec<&str>>()
        .try_into()
        .ok()
        .ok_or(ParseError::ParseStringError(s.into()))
}
