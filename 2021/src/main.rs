use std::collections::HashMap;
use std::env;
use std::error::Error;

mod helpers;

macro_rules! mods {
    ( $( $x:ident ),* ) => {
        $(
            mod $x;
        )*

        fn register(registry: &mut Registry) {
        $(
            $x::register(registry);
        )*
        }
    };
}

mods!(
    day1, day2, day3, day4, day5, day6, day7, day8, day9, day10, day11, day12, day13, day14, day15
);

type EmptyResult = Result<(), Box<dyn Error>>;
type Func = fn() -> EmptyResult;

pub struct Registry {
    impls: HashMap<String, Func>,
}

impl Registry {
    fn new() -> Self {
        Registry {
            impls: HashMap::new(),
        }
    }
    pub fn register(self: &mut Self, handle: String, f: Func) {
        self.impls.insert(handle, f);
    }
}

fn main() -> EmptyResult {
    let mut registry = Registry::new();
    register(&mut registry);

    let mut args = env::args();
    args.next();
    match args.next() {
        Some(arg) => {
            if let Some(part) = registry.impls.get(arg.as_str()) {
                part()
            } else {
                Err("no impl".into())
            }
        }
        _ => Err("no arg".into()),
    }
}
