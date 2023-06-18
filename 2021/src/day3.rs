use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::io::{self, BufRead};
use std::{fs, str};

pub fn register(registry: &mut Registry) {
    registry.register("3p1".into(), p1);
    registry.register("3p2".into(), p2);
}

const WIDTH: usize = 12;

#[derive(Debug, Clone)]
struct Bitcounts([i32; WIDTH]);

impl Bitcounts {
    fn gamma(self: &Self) -> u32 {
        self.int()
    }
    fn epsilon(self: &Self) -> u32 {
        self.flip().int()
    }
    fn int(self: &Self) -> u32 {
        self.0.iter().fold(
            0,
            |state, s| if *s > 0 { state << 1 | 1 } else { state << 1 },
        )
    }
    fn flip(&self) -> Self {
        Self(self.0.map(|i| -i))
    }
}

impl std::ops::Add for &Bitcounts {
    type Output = Bitcounts;
    fn add(self, rhs: Self) -> Self::Output {
        self.0
            .iter()
            .zip(rhs.0.iter())
            .map(|(l, r)| *l + *r)
            .collect::<Vec<i32>>()
            .try_into()
            .unwrap()
    }
}

impl<'a> std::iter::Sum<&'a Bitcounts> for Bitcounts {
    fn sum<I: std::iter::Iterator<Item = &'a Bitcounts>>(iter: I) -> Self {
        iter.fold(
            Bitcounts([0; WIDTH]),
            |mut state: Bitcounts, curr: &Bitcounts| -> Bitcounts {
                state = &state + curr;
                state
            },
        )
    }
}

impl TryFrom<Vec<i32>> for Bitcounts {
    type Error = ParseError;
    fn try_from(v: Vec<i32>) -> Result<Self, ParseError> {
        Ok(Bitcounts(TryInto::<[i32; WIDTH]>::try_into(v).map_err(
            |v| ParseError::ParseStringError(format!("{:?}", v)),
        )?))
    }
}

impl std::str::FromStr for Bitcounts {
    type Err = ParseError;
    fn from_str(s: &str) -> Result<Self, Self::Err> {
        Ok(s.as_bytes()
            .iter()
            .map(|c| -> Result<i32, ParseError> {
                match c {
                    b'0' => Ok(-1),
                    b'1' => Ok(1),
                    _ => Err(s.into()),
                }
            })
            .collect::<Result<Vec<_>, ParseError>>()?
            .try_into()?)
    }
}

fn parse_line(line: Result<String, io::Error>) -> Result<Bitcounts, ParseError> {
    line?.parse()
}

fn parse_file() -> Result<impl Iterator<Item = Result<Bitcounts, ParseError>>, ParseError> {
    let f = fs::File::open("./inputs/3.txt")?;
    Ok(io::BufReader::new(f).lines().map(parse_line))
}

fn p1() -> EmptyResult {
    let s = parse_file()?
        .collect::<Result<Vec<_>, ParseError>>()?
        .iter()
        .sum::<Bitcounts>();
    println!("{}", s.gamma() * s.epsilon());
    Ok(())
}

#[derive(Debug, Clone)]
struct BitcountsVec(Vec<Bitcounts>);

impl BitcountsVec {
    fn o2(mut self: Self) -> u32 {
        let mut s = self.0.iter().sum::<Bitcounts>();
        let mut bit = 0;
        while self.0.len() > 1 && bit < WIDTH {
            let want = if s.0[bit] >= 0 { 1 } else { -1 };
            self.0.retain(|e| e.0[bit] == want);
            s = self.0.iter().sum::<Bitcounts>();
            bit += 1;
        }
        if self.0.len() != 1 {
            panic!("filtered len != 1")
        }
        self.0[0].int()
    }
    fn co2(mut self: Self) -> u32 {
        let mut s = self.0.iter().sum::<Bitcounts>();
        let mut bit = 0;
        while self.0.len() > 1 && bit < WIDTH {
            let want = if s.0[bit] >= 0 { -1 } else { 1 };
            self.0.retain(|e| e.0[bit] == want);
            s = self.0.iter().sum::<Bitcounts>();
            bit += 1;
        }
        if self.0.len() != 1 {
            panic!("filtered len != 1")
        }
        self.0[0].int()
    }
}

fn p2() -> EmptyResult {
    let bitcounts = BitcountsVec(parse_file()?.collect::<Result<Vec<_>, ParseError>>()?);
    println!("{:?}", bitcounts.clone().o2() * bitcounts.co2());
    Ok(())
}
