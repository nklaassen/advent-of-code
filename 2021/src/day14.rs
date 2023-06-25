use crate::helpers::ParseError;
use crate::{EmptyResult, Registry};
use std::collections;
use std::fs;

pub fn register(registry: &mut Registry) {
    registry.register("14p1".into(), p1);
    registry.register("14p2".into(), p2);
}

#[derive(Debug)]
struct Input {
    template: String,
    rules: Vec<Rule>,
}

#[derive(Debug)]
struct Rule {
    pair: (char, char),
    insert: char,
}

fn parse_file() -> Result<Input, ParseError> {
    let f = fs::read_to_string("./inputs/14.txt")?;
    let chunks = f.split_once("\n\n").unwrap();
    let template = chunks.0.to_owned();
    let rules: Vec<Rule> = chunks
        .1
        .lines()
        .map(|line| {
            let chunks = line.split_once(" -> ").unwrap();
            let mut pair_chars = chunks.0.chars();
            let pair = (pair_chars.next().unwrap(), pair_chars.next().unwrap());
            let insert = chunks.1.chars().next().unwrap();
            Rule { pair, insert }
        })
        .collect();
    Ok(Input { template, rules })
}

fn apply(template: &String, rules: &Vec<Rule>) -> String {
    let mut chain: String = template
        .chars()
        .zip(template[1..template.len()].chars())
        .flat_map(|pair| {
            for rule in rules {
                if rule.pair == pair {
                    return vec![pair.0, rule.insert].into_iter();
                }
            }
            vec![pair.0].into_iter()
        })
        .collect();
    chain.push(template.chars().last().unwrap());
    chain
}

fn p1() -> EmptyResult {
    let input = parse_file()?;
    let mut chain = input.template;
    //println!("0: {}", chain);
    for _i in 1..=10 {
        chain = apply(&chain, &input.rules);
        //println!("{}: {}", _i, chain);
    }
    let mut counts: collections::HashMap<char, u32> = collections::HashMap::new();
    chain.chars().for_each(|c| {
        counts.entry(c).and_modify(|e| *e += 1).or_insert(1);
    });
    let max_count = counts.values().max().unwrap();
    let min_count = counts.values().min().unwrap();
    println!("{}", max_count - min_count);
    Ok(())
}

fn apply2(
    pair_counts: &mut collections::HashMap<(char, char), i64>,
    rules: &collections::HashMap<(char, char), [(char, char); 2]>,
) {
    let mut deltas: collections::HashMap<(char, char), i64> = collections::HashMap::new();
    pair_counts
        .iter()
        .for_each(|(pair, &count)| match rules.get(pair) {
            None => {}
            Some(rule) => {
                deltas
                    .entry(*pair)
                    .and_modify(|delta| *delta -= count)
                    .or_insert(-count);
                for new_pair in rule {
                    deltas
                        .entry(*new_pair)
                        .and_modify(|delta| *delta += count)
                        .or_insert(count);
                }
            }
        });
    deltas.iter().for_each(|(pair, &count)| {
        pair_counts
            .entry(*pair)
            .and_modify(|c| *c += count)
            .or_insert(count);
    });
}

fn p2() -> EmptyResult {
    let input = parse_file()?;
    let mut pair_counts: collections::HashMap<(char, char), i64> = input
        .template
        .chars()
        .zip(input.template[1..input.template.len()].chars())
        .fold(collections::HashMap::new(), |mut counts, pair| {
            counts.entry(pair).and_modify(|c| *c += 1).or_insert(1);
            counts
        });
    let rules: collections::HashMap<(char, char), [(char, char); 2]> =
        input
            .rules
            .iter()
            .fold(collections::HashMap::new(), |mut rules, rule| {
                rules.insert(
                    rule.pair,
                    [(rule.pair.0, rule.insert), (rule.insert, rule.pair.1)],
                );
                rules
            });
    //println!("step 0: {:?}", pair_counts);
    for _i in 1..=40 {
        apply2(&mut pair_counts, &rules);
        //println!("step {}: {:?}", _i, pair_counts);
    }
    let mut char_counts: collections::HashMap<char, i64> = collections::HashMap::new();
    char_counts.insert(input.template.chars().next().unwrap(), 1);
    pair_counts.iter().for_each(|(pair, &count)| {
        char_counts
            .entry(pair.1)
            .and_modify(|c| *c += count)
            .or_insert(count);
    });
    //println!("char counts: {:?}", char_counts);
    let max_count = char_counts.values().max().unwrap();
    let min_count = char_counts.values().min().unwrap();
    println!("{}", max_count - min_count);
    Ok(())
}
