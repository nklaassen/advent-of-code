use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::collections::HashSet;
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("9p1".into(), p1);
    registry.register("9p2".into(), p2);
}

fn parse_line(line: Result<String, io::Error>) -> Result<Vec<u8>, ParseError> {
    Ok(line?
        .chars()
        .map(|c| c.to_digit(10).unwrap() as u8)
        .collect())
}

fn parse_file() -> Result<Vec<Vec<u8>>, ParseError> {
    let f = fs::File::open("./inputs/9.txt")?;
    io::BufReader::new(f)
        .lines()
        .map(parse_line)
        .collect::<Result<Vec<_>, _>>()
}

const DIRECTIONS: [(isize, isize); 4] = [(0, 1), (1, 0), (0, -1), (-1, 0)];

fn adjacents(
    max_i: usize,
    max_j: usize,
    i: usize,
    j: usize,
) -> impl Iterator<Item = (usize, usize)> {
    let max_i = max_i as isize;
    let max_j = max_j as isize;
    let i = i as isize;
    let j = j as isize;
    DIRECTIONS
        .map(|d| (i as isize + d.0, j as isize + d.1))
        .into_iter()
        .filter_map(move |p| {
            if p.0 >= 0 && p.0 < max_i && p.1 >= 0 && p.1 < max_j {
                Some((p.0 as usize, p.1 as usize))
            } else {
                None
            }
        })
}

fn p1() -> EmptyResult {
    let map = parse_file()?;

    let max_i = map.len();
    let max_j = map[0].len();
    let mut risk_sum: u32 = 0;
    for (i, row) in map.iter().enumerate() {
        for (j, &height) in row.iter().enumerate() {
            let mut low = true;
            for (ai, aj) in adjacents(max_i, max_j, i, j) {
                if map[ai][aj] <= height {
                    low = false;
                    break;
                }
            }
            if low {
                risk_sum += height as u32 + 1;
            }
        }
    }
    println!("{}", risk_sum);
    Ok(())
}

fn p2() -> EmptyResult {
    let map = parse_file()?;

    let max_i = map.len();
    let max_j = map[0].len();

    let adjacents = |(i, j)| adjacents(max_i, max_j, i, j);

    let mut low_points: Vec<(usize, usize)> = Vec::new();
    for (i, row) in map.iter().enumerate() {
        for (j, &height) in row.iter().enumerate() {
            let mut low = true;
            for (ai, aj) in adjacents((i, j)) {
                if map[ai][aj] <= height {
                    low = false;
                    break;
                }
            }
            if low {
                low_points.push((i, j));
            }
        }
    }

    let mut basin_sizes: Vec<usize> = Vec::new();
    for &low_point in &low_points {
        let mut seen: HashSet<(usize, usize)> = HashSet::new();
        seen.insert(low_point);
        let mut pending: Vec<(usize, usize)> = vec![low_point];
        while let Some(curr) = pending.pop() {
            let curr_height = map[curr.0][curr.1];
            for a in adjacents(curr) {
                if seen.contains(&a) {
                    continue;
                }
                let adjacent_height = map[a.0][a.1];
                if adjacent_height < 9 && adjacent_height >= curr_height {
                    seen.insert(a);
                    pending.push(a);
                }
            }
        }
        basin_sizes.push(seen.len());
    }
    basin_sizes.sort_unstable();
    let ans = basin_sizes[basin_sizes.len() - 3..]
        .iter()
        .product::<usize>();
    println!("{}", ans);
    Ok(())
}
