use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("7p1".into(), p1);
    registry.register("7p2".into(), p2);
}

fn parse_file() -> Result<Vec<u32>, ParseError> {
    let f = fs::File::open("./inputs/7.txt")?;
    io::BufReader::new(f)
        .split(b',')
        .map(|s| -> Result<_, ParseError> {
            String::from_utf8_lossy(s?.as_slice())
                .to_string()
                .trim()
                .parse()
                .map_err(ParseError::ParseIntError)
        })
        .collect::<Result<Vec<_>, ParseError>>()
}

fn p1() -> EmptyResult {
    let input = parse_file()?;

    let max = input.iter().max().unwrap();
    let min = input.iter().min().unwrap();

    let fuel = (*min..=*max)
        .map(|i| input.iter().map(|x| x.abs_diff(i)).sum::<u32>())
        .min()
        .unwrap();

    println!("{:?}", fuel);
    Ok(())
}

fn p2() -> EmptyResult {
    let input = parse_file()?;

    let max = input.iter().max().unwrap();
    let min = input.iter().min().unwrap();

    let fuel = (*min..=*max)
        .map(|i| {
            input
                .iter()
                .map(|x| -> u32 {
                    let d = x.abs_diff(i);
                    d * (d + 1) / 2
                })
                .sum::<u32>()
        })
        .min()
        .unwrap();

    println!("{:?}", fuel);
    Ok(())
}
