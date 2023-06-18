use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("5p1".into(), p1);
    registry.register("5p2".into(), p2);
}

type Point = [i32; 2];
type Line = [Point; 2];

fn parse_line(line: Result<String, io::Error>) -> Result<Line, ParseError> {
    line?
        .split(" -> ")
        .map(|p| -> Result<Point, ParseError> {
            p.split(",")
                .map(|coord| coord.parse::<i32>())
                .collect::<Result<Vec<i32>, std::num::ParseIntError>>()?
                .try_into()
                .map_err(|_| ParseError::ParseStringError("incorrect point length".into()))
        })
        .collect::<Result<Vec<_>, _>>()?
        .try_into()
        .map_err(|_| ParseError::ParseStringError("incorrect line length".into()))
}

fn parse_file() -> Result<impl Iterator<Item = Result<Line, ParseError>>, ParseError> {
    let f = fs::File::open("./inputs/5.txt")?;
    Ok(io::BufReader::new(f).lines().map(parse_line))
}

fn add(p1: Point, p2: Point) -> Point {
    [p1[0] + p2[0], p1[1] + p2[1]]
}

fn sub(p1: Point, p2: Point) -> Point {
    [p1[0] - p2[0], p1[1] - p2[1]]
}

fn scale_inverse(p: Point, s: i32) -> Point {
    [p[0] / s, p[1] / s]
}

fn delta(p1: Point, p2: Point) -> Point {
    let diff = sub(p2, p1);
    let max = diff.iter().map(|x| x.abs()).max().unwrap();
    scale_inverse(diff, max)
}

fn p1() -> EmptyResult {
    let lines = parse_file()?.collect::<Result<Vec<_>, _>>()?;
    let lines = lines
        .into_iter()
        .filter(|line| (0..2).any(|i| line[0][i] == line[1][i]))
        .collect::<Vec<Line>>();
    let mut field: std::collections::HashMap<Point, u32> = std::collections::HashMap::new();

    for line in &lines {
        let d = delta(line[0], line[1]);
        let mut curr: Point = line[0];
        *field.entry(curr).or_default() += 1;
        while {
            curr = add(curr, d);
            *field.entry(curr).or_default() += 1;
            curr != line[1]
        } {}
    }

    println!("{}", field.into_values().filter(|&v| v > 1).count());
    Ok(())
}

fn p2() -> EmptyResult {
    let lines = parse_file()?.collect::<Result<Vec<_>, _>>()?;
    let mut field: std::collections::HashMap<Point, u32> = std::collections::HashMap::new();

    for line in &lines {
        let d = delta(line[0], line[1]);
        let mut curr: Point = line[0];
        *field.entry(curr).or_default() += 1;
        while {
            curr = add(curr, d);
            *field.entry(curr).or_default() += 1;
            curr != line[1]
        } {}
    }

    println!("{}", field.into_values().filter(|&v| v > 1).count());
    Ok(())
}
