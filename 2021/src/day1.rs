use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("1p1".into(), p1);
    registry.register("1p2".into(), p2);
}

fn parse_line(line: Result<String, io::Error>) -> Result<i64, ParseError> {
    Ok(line?.parse()?)
}

fn parse_file(file: fs::File) -> impl Iterator<Item = Result<i64, ParseError>> {
    io::BufReader::new(file).lines().map(parse_line)
}

fn p1() -> EmptyResult {
    let file = fs::File::open("./inputs/1.txt")?;
    println!(
        "{}",
        parse_file(file)
            .fold(
                Ok((std::i64::MAX, 0)),
                |prev, curr| -> Result<(i64, i32), ParseError> {
                    let prev = prev?;
                    let curr = curr?;
                    if curr > prev.0 {
                        Ok((curr, prev.1 + 1))
                    } else {
                        Ok((curr, prev.1))
                    }
                }
            )?
            .1
    );
    Ok(())
}

fn p2() -> EmptyResult {
    let file = fs::File::open("./inputs/1.txt")?;
    println!(
        "{}",
        parse_file(file)
            .collect::<Result<Vec<i64>, ParseError>>()?
            .windows(3)
            .fold((std::i64::MAX, 0), |prev, window| {
                let s = window.iter().sum();
                if s > prev.0 {
                    (s, prev.1 + 1)
                } else {
                    (s, prev.1)
                }
            })
            .1,
    );
    Ok(())
}
