use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("6p1".into(), p1);
    registry.register("6p2".into(), p2);
}

type Input = Vec<u8>;

fn parse_file() -> Result<Input, ParseError> {
    let f = fs::File::open("./inputs/6.txt")?;
    io::BufReader::new(f)
        .split(b',')
        .map(|s| -> Result<u8, ParseError> {
            String::from_utf8_lossy(s?.as_slice())
                .to_string()
                .trim()
                .parse()
                .map_err(ParseError::ParseIntError)
        })
        .collect::<Result<Vec<_>, ParseError>>()
}

const P1_DAYS: u8 = 80;
const NEWBORN_DAYS: u8 = 8;
const RESET_DAYS: u8 = 6;

fn p1() -> EmptyResult {
    let mut state = parse_file()?;
    for i in 0..P1_DAYS {
        let mut new_fish = 0;
        for fish in &mut state {
            if *fish == 0 {
                *fish = RESET_DAYS;
                new_fish += 1;
            } else {
                *fish -= 1;
            }
        }
        state.append(&mut vec![NEWBORN_DAYS; new_fish]);
        println!("{}", i);
    }
    println!("{}", state.len());
    Ok(())
}

const P2_DAYS: u16 = 256;

fn p2() -> EmptyResult {
    let input = parse_file()?;
    // count of fish at each number of remaining days
    let mut state = [0; NEWBORN_DAYS as usize + 1];
    for days in &input {
        state[*days as usize] += 1
    }
    for _ in 0..P2_DAYS {
        state.rotate_left(1);
        state[RESET_DAYS as usize] += state[NEWBORN_DAYS as usize];
    }
    println!("{}", state.iter().sum::<u64>());
    Ok(())
}
