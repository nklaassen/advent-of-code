use crate::helpers::{split_fields, ParseError};
use crate::{EmptyResult, Registry};
use std::io::{self, BufRead};
use std::{fs, str};

pub fn register(registry: &mut Registry) {
    registry.register("2p1".into(), p1);
    registry.register("2p2".into(), p2);
}

#[derive(Debug)]
enum Direction {
    Forward,
    Down,
    Up,
}

impl str::FromStr for Direction {
    type Err = ParseError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "forward" => Ok(Direction::Forward),
            "up" => Ok(Direction::Up),
            "down" => Ok(Direction::Down),
            _ => Err(ParseError::ParseStringError(s.into())),
        }
    }
}

#[derive(Debug)]
struct Command {
    dir: Direction,
    delta: i32,
}

impl str::FromStr for Command {
    type Err = ParseError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let fields: [&str; 2] = split_fields(s)?;
        Ok(Command {
            dir: fields[0].parse()?,
            delta: fields[1].parse()?,
        })
    }
}

fn parse_line(line: Result<String, io::Error>) -> Result<Command, ParseError> {
    line?.parse()
}

fn parse_file() -> Result<impl Iterator<Item = Result<Command, ParseError>>, ParseError> {
    let f = fs::File::open("./inputs/2.txt")?;
    Ok(io::BufReader::new(f).lines().map(parse_line))
}

#[derive(Debug)]
struct State {
    forward: i32,
    depth: i32,
}

fn p1() -> EmptyResult {
    let final_state = parse_file()?.fold(
        Ok(State {
            forward: 0,
            depth: 0,
        }),
        |s: Result<State, ParseError>, c| {
            let s = s?;
            let c = c?;
            Ok(match c.dir {
                Direction::Forward => State {
                    forward: s.forward + c.delta,
                    depth: s.depth,
                },
                Direction::Down => State {
                    forward: s.forward,
                    depth: s.depth + c.delta,
                },
                Direction::Up => State {
                    forward: s.forward,
                    depth: s.depth - c.delta,
                },
            })
        },
    )?;
    println!("{}", final_state.depth * final_state.forward);
    Ok(())
}

#[derive(Debug)]
struct State2 {
    forward: i32,
    depth: i32,
    aim: i32,
}

fn p2() -> EmptyResult {
    let final_state = parse_file()?.fold(
        Ok(State2 {
            forward: 0,
            depth: 0,
            aim: 0,
        }),
        |s: Result<State2, ParseError>, c| {
            let s = s?;
            let c = c?;
            Ok(match c.dir {
                Direction::Forward => State2 {
                    forward: s.forward + c.delta,
                    depth: s.depth + s.aim * c.delta,
                    aim: s.aim,
                },
                Direction::Down => State2 {
                    forward: s.forward,
                    depth: s.depth,
                    aim: s.aim + c.delta,
                },
                Direction::Up => State2 {
                    forward: s.forward,
                    depth: s.depth,
                    aim: s.aim - c.delta,
                },
            })
        },
    )?;
    println!("{:?}", final_state);
    println!("{}", final_state.depth * final_state.forward);
    Ok(())
}
