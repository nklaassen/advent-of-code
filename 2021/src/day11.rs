use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("11p1".into(), p1);
    registry.register("11p2".into(), p2);
}

const WIDTH: usize = 10;
const HEIGHT: usize = 10;

fn parse_file() -> Result<[u8; WIDTH * HEIGHT], ParseError> {
    let f = fs::File::open("./inputs/11.txt")?;
    io::BufReader::new(f)
        .lines()
        .map(|line| -> Result<_, ParseError> { Ok(line?.into_bytes()) })
        .collect::<Result<Vec<Vec<u8>>, ParseError>>()?
        .into_iter()
        .fold(Vec::new(), |mut state, mut curr| {
            state.append(&mut curr);
            state
        })
        .iter()
        .map(|b| b - b'0')
        .collect::<Vec<u8>>()
        .try_into()
        .map_err(|_| ParseError::ParseStringError("bad input size".into()))
}

fn step(octopi: &mut [u8; 100], flashes: &mut u32) -> bool {
    for o in octopi.iter_mut() {
        *o += 1
    }
    for i in 0..octopi.len() {
        if octopi[i] == 10 {
            flash(flashes, octopi, i);
        }
    }
    let mut all_flashed = true;
    for o in octopi {
        if *o > 9 {
            *o = 0
        } else {
            all_flashed = false
        }
    }
    all_flashed
}
fn flash(flashes: &mut u32, octopi: &mut [u8; 100], i: usize) {
    *flashes += 1;
    octopi[i] = 11;
    for a in adjacents(i) {
        match octopi[a] {
            x if x < 9 => octopi[a] += 1,
            x if x == 9 => flash(flashes, octopi, a),
            _ => {}
        }
    }
}
const OFFSETS: [isize; 8] = [-11, -10, -9, -1, 1, 9, 10, 11];
fn adjacents(i: usize) -> impl Iterator<Item = usize> {
    OFFSETS.iter().flat_map(move |o| match i as isize + o {
        a if a >= 0
            && a < (WIDTH * HEIGHT) as isize
            && (i % 10 > 0 || a % 10 < 2)
            && (i % 10 < 9 || a % 10 > 7) =>
        {
            Some(a as usize)
        }
        _ => None,
    })
}

fn print(v: [u8; WIDTH * HEIGHT]) -> EmptyResult {
    let mut s: Vec<u8> = Vec::new();
    for j in 0..HEIGHT {
        for i in 0..WIDTH {
            s.push(v[j * WIDTH + i] + b'0');
        }
        s.push(b'\n');
    }
    println!("{}", String::from_utf8(s)?);
    Ok(())
}

fn p1() -> EmptyResult {
    let mut octopi = parse_file()?;
    let mut flashes = 0;
    for _ in 0..100 {
        step(&mut octopi, &mut flashes);
    }
    print(octopi)?;
    println!("{}", flashes);
    Ok(())
}

fn p2() -> EmptyResult {
    let mut octopi = parse_file()?;
    let mut _flashes = 0;
    let mut i = 0;
    loop {
        if step(&mut octopi, &mut _flashes) {
            break;
        }
        i += 1
    }
    println!("{}", i + 1);
    Ok(())
}
