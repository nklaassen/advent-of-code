use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::collections;
use std::fs;

pub fn register(registry: &mut Registry) {
    registry.register("13p1".into(), p1);
    registry.register("13p2".into(), p2);
}

#[derive(Debug)]
struct Input {
    points: Vec<Point>,
    folds: Vec<Fold>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
struct Point {
    x: i32,
    y: i32,
}

#[derive(Debug, Clone, Copy)]
struct Fold {
    axis: Axis,
    line: i32,
}

#[derive(Debug, Clone, Copy, PartialEq)]
enum Axis {
    X,
    Y,
}

fn parse_file() -> Result<Input, ParseError> {
    let f = fs::read_to_string("./inputs/13.txt")?;
    let chunks = f.split_once("\n\n").unwrap();
    let points = chunks.0;
    let points: Vec<Point> = points
        .lines()
        .map(|line| {
            line.split_once(",")
                .map(|p| Point {
                    x: p.0.parse().unwrap(),
                    y: p.1.parse().unwrap(),
                })
                .unwrap()
        })
        .collect();
    let folds = chunks.1;
    let folds: Vec<Fold> = folds
        .lines()
        .map(|line| line.strip_prefix("fold along ").unwrap())
        .map(|line| {
            line.split_once("=")
                .map(|p| Fold {
                    axis: if p.0 == "x" { Axis::X } else { Axis::Y },
                    line: p.1.parse().unwrap(),
                })
                .unwrap()
        })
        .collect();
    Ok(Input { points, folds })
}

fn fold_up(p: &Point, y: i32) -> Point {
    Point {
        x: p.x,
        y: if p.y > y { y - (p.y - y) } else { p.y },
    }
}

fn fold_left(p: &Point, x: i32) -> Point {
    Point {
        y: p.y,
        x: if p.x > x { x - (p.x - x) } else { p.x },
    }
}

fn fold(f: &Fold, p: &Point) -> Point {
    if f.axis == Axis::X {
        fold_left(p, f.line)
    } else {
        fold_up(p, f.line)
    }
}

fn p1() -> EmptyResult {
    let input = parse_file()?;

    let first_fold = input.folds[0];
    let mut final_points: collections::HashSet<Point> = collections::HashSet::new();
    for p in &input.points {
        final_points.insert(fold(&first_fold, p));
    }
    println!("{:?}", final_points.len());
    Ok(())
}

fn print(points: &collections::HashSet<Point>) {
    let max_x = points.iter().map(|p| p.x).max().unwrap();
    let max_y = points.iter().map(|p| p.y).max().unwrap();
    for y in 0..=max_y {
        for x in 0..=max_x {
            if points.contains(&Point { x, y }) {
                print!("x");
            } else {
                print!(".");
            }
        }
        println!();
    }
}

fn p2() -> EmptyResult {
    let input = parse_file()?;
    let first_fold = input.folds[0];
    let mut curr_points: collections::HashSet<Point> = collections::HashSet::new();
    for p in &input.points {
        curr_points.insert(fold(&first_fold, p));
    }
    for f in &input.folds[1..] {
        let mut next_points = collections::HashSet::new();
        for p in curr_points.iter() {
            next_points.insert(fold(f, p));
        }
        curr_points = next_points;
    }
    print(&curr_points);
    Ok(())
}
