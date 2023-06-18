use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("tp1".into(), p1);
    registry.register("tp2".into(), p2);
}

fn parse_line(line: Result<String, io::Error>) -> Result<Command, ParseError> {
    line?.parse()
}

fn parse_file() -> Result<impl Iterator<Item = Result<Command, ParseError>>, ParseError> {
    let f = fs::File::open("./inputs/t.txt")?;
    Ok(io::BufReader::new(f).lines().map(parse_line))
}

fn p1() -> EmptyResult {
    todo!()
}

fn p2() -> EmptyResult {
    todo!()
}
