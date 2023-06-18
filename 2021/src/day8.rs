use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("8p1".into(), p1);
    registry.register("8p2".into(), p2);
}

const PATTERNS_PER_ENTRY: usize = 10;
const OUTPUTS_PER_ENTRY: usize = 4;

const UNIQUES: [(u32, usize); 4] = [(1, 2), (4, 4), (7, 3), (8, 7)];

#[derive(Debug)]
struct Entry {
    patterns: [String; PATTERNS_PER_ENTRY],
    outputs: [String; OUTPUTS_PER_ENTRY],
}

impl std::str::FromStr for Entry {
    type Err = ParseError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let chunks = s.split(" | ");
        let mut chunks = chunks.map(|chunk| {
            chunk
                .split_whitespace()
                .map(|s| s.to_owned())
                .collect::<Vec<String>>()
        });
        let patterns = chunks
            .next()
            .unwrap()
            .try_into()
            .map_err(|_| ParseError::ParseStringError("wrong pattern count".into()))?;
        let outputs = chunks
            .next()
            .unwrap()
            .try_into()
            .map_err(|_| ParseError::ParseStringError("wrong output count".into()))?;
        Ok(Entry { patterns, outputs })
    }
}

fn parse_line(line: Result<String, io::Error>) -> Result<Entry, ParseError> {
    line?.parse()
}

fn parse_file() -> Result<impl Iterator<Item = Result<Entry, ParseError>>, ParseError> {
    let f = fs::File::open("./inputs/8.txt")?;
    Ok(io::BufReader::new(f).lines().map(parse_line))
}

fn p1() -> EmptyResult {
    let entries = parse_file()?.collect::<Result<Vec<_>, _>>()?;

    let unique_lengths: Vec<usize> = UNIQUES.iter().map(|u| u.1).collect();

    let unique_output_count = entries
        .iter()
        .flat_map(|e| &e.outputs)
        .filter(|output| unique_lengths.contains(&output.len()))
        .count();
    println!("{}", unique_output_count);

    Ok(())
}

fn digit_for_output(patterns: &[String; PATTERNS_PER_ENTRY], output: &String) -> u32 {
    let uniques: std::collections::HashMap<u32, &String> = UNIQUES
        .map(|u| (u.0, patterns.iter().find(|p| p.len() == u.1).unwrap()))
        .into_iter()
        .collect();

    if let Some(u) = UNIQUES.iter().find(|u| u.1 == output.len()) {
        return u.0;
    }

    let one: &String = uniques[&1];
    let four: &String = uniques[&4];

    // only 6 has len 6 and doesn't contain all segments from 1
    if output.len() == 6 && !one.chars().all(|c| output.contains(c)) {
        return 6;
    }

    // only 0 has len 6 and doesn't contain all segments from 4
    if output.len() == 6 && !four.chars().all(|c| output.contains(c)) {
        return 0;
    }

    // only remaining with len 6 is 9
    if output.len() == 6 {
        return 9;
    }

    // only 3 has len 5 and all segments from 1
    if output.len() == 5 && one.chars().all(|c| output.contains(c)) {
        return 3;
    }

    // only remaining with 2 common segments with 4 is 2
    if four.chars().filter(|&c| output.contains(c)).count() == 2 {
        return 2;
    }

    5
}

fn p2() -> EmptyResult {
    let mut sum = 0;
    for entry in parse_file()? {
        let entry = entry?;
        let output: u32 = entry
            .outputs
            .iter()
            .map(|o| digit_for_output(&entry.patterns, &o))
            .into_iter()
            .reduce(|prev, curr| prev * 10 + curr)
            .unwrap();
        sum += output;
    }
    println!("{}", sum);
    Ok(())
}
