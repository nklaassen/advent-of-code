use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

type Map = std::collections::HashMap<char, char>;

pub fn register(registry: &mut Registry) {
    registry.register("10p1".into(), p1);
    registry.register("10p2".into(), p2);
}

fn parse_file() -> io::Result<Vec<String>> {
    let f = fs::File::open("./inputs/10.txt")?;
    io::BufReader::new(f)
        .lines()
        .collect::<Result<Vec<String>, _>>()
}

#[derive(Debug)]
enum ParseResult {
    Valid,
    Incomplete(Vec<char>),
    Corrupted(char),
}

fn parse(closers: &Map, s: &String) -> ParseResult {
    let mut stack: Vec<char> = Vec::new();
    for c in s.chars() {
        if let Some(closer) = closers.get(&c) {
            stack.push(*closer);
            continue;
        }
        match stack.pop() {
            None => return ParseResult::Corrupted(c),
            Some(closer) => {
                if closer == c {
                    continue;
                }
                return ParseResult::Corrupted(c);
            }
        }
    }
    if stack.len() > 0 {
        stack.reverse();
        return ParseResult::Incomplete(stack);
    }
    ParseResult::Valid
}

fn p1() -> EmptyResult {
    let closers: Map = Map::from([('(', ')'), ('[', ']'), ('{', '}'), ('<', '>')]);
    let lines = parse_file()?;
    let score: u32 = lines
        .iter()
        .filter_map(|line| match parse(&closers, line) {
            ParseResult::Corrupted(c) => Some(c),
            _ => None,
        })
        .map(|c| match c {
            ')' => 3,
            ']' => 57,
            '}' => 1197,
            '>' => 25137,
            _ => panic!("bad input"),
        })
        .sum();
    println!("{}", score);
    Ok(())
}

fn p2() -> EmptyResult {
    let closers: Map = Map::from([('(', ')'), ('[', ']'), ('{', '}'), ('<', '>')]);
    let lines = parse_file()?;
    let mut scores: Vec<u64> = lines
        .iter()
        .filter_map(|line| match parse(&closers, line) {
            ParseResult::Incomplete(v) => Some(v),
            _ => None,
        })
        .map(|s| {
            s.iter()
                .map(|c| match c {
                    ')' => 1,
                    ']' => 2,
                    '}' => 3,
                    '>' => 4,
                    _ => panic!("bad input"),
                })
                .fold(0_u64, |state, n| state * 5 + n)
        })
        .collect();
    scores.sort_unstable();
    let mid = scores[scores.len() / 2];
    println!("{}", mid);
    Ok(())
}
