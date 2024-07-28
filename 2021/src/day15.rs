use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("15p1".into(), p1);
    registry.register("15p2".into(), p2);
}

type RiskMap = Vec<Vec<u32>>;

fn parse_file() -> Result<RiskMap, ParseError> {
    let f = fs::File::open("./inputs/15.txt")?;
    io::BufReader::new(f)
        .lines()
        .map(|line| Ok(line?.bytes().map(|b| (b - b'0') as u32).collect()))
        .collect()
}

#[derive(Copy, Clone, Eq, PartialEq, Hash, Debug)]
struct Point {
    x: usize,
    y: usize,
}

#[derive(Copy, Clone, Eq, PartialEq, Debug)]
struct State {
    p: Point,
    cost: u32,
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        other.cost.cmp(&self.cost)
    }
}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

fn p1() -> EmptyResult {
    let risk_map = parse_file()?;
    println!("{:?}", risk_map);

    let mut seen = std::collections::HashSet::<Point>::new();
    let mut costs = std::collections::BinaryHeap::<State>::new();
    costs.push(State {
        p: Point { x: 0, y: 0 },
        cost: 0,
    });

    let target = Point {
        x: risk_map.len() - 1,
        y: risk_map[0].len() - 1,
    };
    println!("target: {:?}", target);
    while let Some(curr) = costs.pop() {
        if seen.contains(&curr.p) {
            continue;
        }
        seen.insert(curr.p);
        println!("curr: {:?}", curr);
        if curr.p == target {
            println!("{}", curr.cost);
            return Ok(());
        }
        for d in [(-1, 0), (1, 0), (0, -1), (0, 1)].iter() {
            if curr.p.x == 0 && d.0 == -1
                || curr.p.x == risk_map.len() - 1 && d.0 == 1
                || curr.p.y == 0 && d.1 == -1
                || curr.p.y == risk_map[0].len() - 1 && d.1 == 1
            {
                continue;
            }
            let np = Point {
                x: (curr.p.x as i32 + d.0) as usize,
                y: (curr.p.y as i32 + d.1) as usize,
            };
            let nc = curr.cost + risk_map[np.x][np.y];
            costs.push(State { p: np, cost: nc });
        }
    }
    todo!()
}

fn p2() -> EmptyResult {
    let risk_map = parse_file()?;
    //println!("{:?}", risk_map);

    let mut seen = std::collections::HashSet::<Point>::new();
    let mut costs = std::collections::BinaryHeap::<State>::new();
    costs.push(State {
        p: Point { x: 0, y: 0 },
        cost: 0,
    });

    let target = Point {
        x: risk_map.len() * 5 - 1,
        y: risk_map[0].len() * 5 - 1,
    };
    println!("target: {:?}", target);
    while let Some(curr) = costs.pop() {
        if seen.contains(&curr.p) {
            continue;
        }
        seen.insert(curr.p);
        //println!("curr: {:?}", curr);
        if curr.p == target {
            println!("{}", curr.cost);
            return Ok(());
        }
        for d in [(-1, 0), (1, 0), (0, -1), (0, 1)].iter() {
            if curr.p.x == 0 && d.0 == -1
                || curr.p.x == risk_map.len() * 5 - 1 && d.0 == 1
                || curr.p.y == 0 && d.1 == -1
                || curr.p.y == risk_map[0].len() * 5 - 1 && d.1 == 1
            {
                continue;
            }
            let np = Point {
                x: (curr.p.x as i32 + d.0) as usize,
                y: (curr.p.y as i32 + d.1) as usize,
            };
            if seen.contains(&np) {
                continue;
            }
            let mut c = risk_map[np.x % risk_map.len()][np.y % risk_map[0].len()]
                + (np.x / risk_map.len()) as u32
                + (np.y / risk_map[0].len()) as u32;
            if c > 9 {
                c -= 9;
            }
            let nc = curr.cost + c;
            costs.push(State { p: np, cost: nc });
        }
    }
    todo!()
}
