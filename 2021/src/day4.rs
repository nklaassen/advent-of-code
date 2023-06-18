use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;

pub fn register(registry: &mut Registry) {
    registry.register("4p1".into(), p1);
    registry.register("4p2".into(), p2);
}

const WIDTH: usize = 5;

#[derive(Debug)]
struct Board {
    rows: [[u32; WIDTH]; WIDTH],
}

impl std::str::FromStr for Board {
    type Err = ParseError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let rows = s
            .lines()
            .map(|line| -> Result<[u32; WIDTH], ParseError> {
                line.split_whitespace()
                    .map(|s| s.parse())
                    .collect::<Result<Vec<_>, _>>()?
                    .try_into()
                    .map_err(|_| ParseError::ParseStringError("wrong board width".into()))
            })
            .collect::<Result<Vec<_>, _>>()?
            .try_into()
            .map_err(|_| ParseError::ParseStringError("wrong board height".into()))?;
        Ok(Board { rows })
    }
}

impl Board {
    fn win_condition(&self, draws: &Vec<u32>) -> Option<WinCondition> {
        let mut marked = [[false; WIDTH]; WIDTH];
        for (draw_index, draw) in draws.iter().enumerate() {
            for (i, row) in self.rows.iter().enumerate() {
                for (j, value) in row.iter().enumerate() {
                    if *draw == *value {
                        marked[i][j] = true;
                        if Self::winning(marked) {
                            return Some(WinCondition {
                                winning_draw_index: draw_index,
                                winning_draw: *draw,
                                sum_of_unmarked: self.sum_of_unmarked(marked),
                            });
                        }
                    }
                }
            }
        }
        None
    }
    fn winning(marked: [[bool; WIDTH]; WIDTH]) -> bool {
        for row in &marked {
            if row.iter().all(|marked| *marked) {
                return true;
            }
        }
        for j in 0..WIDTH {
            if marked.iter().map(|row| row[j]).all(|marked| marked) {
                return true;
            }
        }
        return false;
    }
    fn sum_of_unmarked(&self, marked: [[bool; WIDTH]; WIDTH]) -> u32 {
        self.rows
            .iter()
            .enumerate()
            .map(move |(i, row)| {
                row.iter()
                    .enumerate()
                    .map(move |(j, value)| ((i, j), value))
            })
            .flatten()
            .filter(|((i, j), _)| !marked[*i][*j])
            .map(|(_, value)| *value)
            .sum()
    }
}

#[derive(Debug)]
struct WinCondition {
    winning_draw_index: usize,
    winning_draw: u32,
    sum_of_unmarked: u32,
}

impl WinCondition {
    fn score(&self) -> u32 {
        self.winning_draw * self.sum_of_unmarked
    }
}

#[derive(Debug)]
struct Input {
    draws: Vec<u32>,
    boards: Vec<Board>,
}

impl std::str::FromStr for Input {
    type Err = ParseError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let mut chunks = s.split("\n\n");
        let draws = chunks
            .next()
            .ok_or(ParseError::ParseStringError(s.into()))?
            .split(",")
            .map(|s| s.parse())
            .collect::<Result<Vec<_>, _>>()?;
        let boards = chunks.map(|s| s.parse()).collect::<Result<Vec<_>, _>>()?;
        Ok(Input { draws, boards })
    }
}

fn parse_file() -> Result<Input, ParseError> {
    fs::read_to_string("./inputs/4.txt")?.parse()
}

fn p1() -> EmptyResult {
    let input = parse_file()?;
    let score: u32 = input
        .boards
        .iter()
        .map(|board| board.win_condition(&input.draws))
        .flatten() // Filter out None
        .min_by_key(|win_condition| win_condition.winning_draw_index)
        .ok_or("no winner")?
        .score();
    println!("{}", score);
    Ok(())
}

fn p2() -> EmptyResult {
    let input = parse_file()?;
    let score: u32 = input
        .boards
        .iter()
        .map(|board| board.win_condition(&input.draws))
        .flatten() // Filter out None
        .max_by_key(|win_condition| win_condition.winning_draw_index)
        .ok_or("no winner")?
        .score();
    println!("{}", score);
    Ok(())
}
