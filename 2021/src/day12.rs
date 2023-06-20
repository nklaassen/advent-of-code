use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::collections;
use std::fs;
use std::io::{self, BufRead};

pub fn register(registry: &mut Registry) {
    registry.register("12p1".into(), p1);
    registry.register("12p2".into(), p2);
}

#[derive(Debug)]
struct Edge {
    a: String,
    b: String,
}

fn parse_line(line: Result<String, io::Error>) -> Result<Edge, ParseError> {
    let line = line?;
    let mut splits = line.split("-");
    let a = splits
        .next()
        .ok_or(ParseError::ParseStringError(line.to_owned()))?
        .to_owned();
    let b = splits
        .next()
        .ok_or(ParseError::ParseStringError(line.to_owned()))?
        .to_owned();
    Ok(Edge { a, b })
}

fn parse_file() -> Result<impl Iterator<Item = Result<Edge, ParseError>>, ParseError> {
    let f = fs::File::open("./inputs/12.txt")?;
    Ok(io::BufReader::new(f).lines().map(parse_line))
}

fn walk(outgoing: &collections::HashMap<String, Vec<String>>) -> u32 {
    let mut seen = collections::HashSet::new();
    seen.insert("start".to_owned());
    return dfs(outgoing, &mut seen, "start");
}

fn dfs(
    outgoing: &collections::HashMap<String, Vec<String>>,
    seen: &mut collections::HashSet<String>,
    start: &str,
) -> u32 {
    let mut count = 0;
    for a in &outgoing[start] {
        if is_lower(&a) && seen.contains(a) {
            continue;
        }
        if a == "end" {
            count += 1;
            continue;
        }
        seen.insert(a.to_owned());
        count += dfs(outgoing, seen, a);
        seen.remove(a);
    }
    count
}

fn is_lower(s: &String) -> bool {
    let b = s.bytes().next().unwrap();
    b >= b'a' && b <= b'z'
}

fn p1() -> EmptyResult {
    let edges: Vec<Edge> = parse_file()?.collect::<Result<Vec<_>, _>>()?;

    let mut outgoing: collections::HashMap<String, Vec<String>> = collections::HashMap::new();
    for edge in &edges {
        outgoing
            .entry(edge.a.clone())
            .and_modify(|v| v.push(edge.b.clone()))
            .or_insert(vec![edge.b.clone()]);
        outgoing
            .entry(edge.b.clone())
            .and_modify(|v| v.push(edge.a.clone()))
            .or_insert(vec![edge.a.clone()]);
    }
    println!("{}", walk(&outgoing));
    Ok(())
}

fn walk2(outgoing: &collections::HashMap<String, Vec<String>>) -> u32 {
    let mut seen = collections::HashSet::new();
    return dfs2(outgoing, &mut seen, "start", "start", false);
}

fn dfs2(
    outgoing: &collections::HashMap<String, Vec<String>>,
    seen: &mut collections::HashSet<String>,
    start: &str,
    prefix: &str,
    exception: bool,
) -> u32 {
    let mut count = 0;
    for a in &outgoing[start] {
        if is_lower(&a) && seen.contains(a) {
            if !exception && a != "end" {
                count += dfs2(outgoing, seen, a, &(prefix.to_owned() + "," + a), true);
            }
            continue;
        }
        if a == "end" {
            count += 1;
            //println!("{}", prefix.to_owned() + "," + "end");
            continue;
        }
        seen.insert(a.to_owned());
        count += dfs2(outgoing, seen, a, &(prefix.to_owned() + "," + a), exception);
        seen.remove(a);
    }
    count
}

fn p2() -> EmptyResult {
    let edges: Vec<Edge> = parse_file()?.collect::<Result<Vec<_>, _>>()?;

    let mut outgoing: collections::HashMap<String, Vec<String>> = collections::HashMap::new();
    for edge in &edges {
        if edge.b != "start" {
            outgoing
                .entry(edge.a.clone())
                .and_modify(|v| v.push(edge.b.clone()))
                .or_insert(vec![edge.b.clone()]);
        }
        if edge.a != "start" {
            outgoing
                .entry(edge.b.clone())
                .and_modify(|v| v.push(edge.a.clone()))
                .or_insert(vec![edge.a.clone()]);
        }
    }
    println!("{}", walk2(&outgoing));
    Ok(())
}
