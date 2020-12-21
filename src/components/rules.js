import React from 'react';
import '../css/header.css';

export default function Rules() {
  return (
    <div className='rulesContainer'>
      <div>The goal is to guess the sequence of four colors in 12 or fewer tries. Colors may be repeated. After each guess, you will receive a ranking of "bulls" and "cows" explained below. Four "bulls" implies that you've won!</div>
      <div><h3>Bulls</h3>Each bull implies a guess that is in the correct location AND is the correct color.</div>
      <div><h3>Cows</h3>Each cow implies a guess that is in the incorrect location BUT is a correct color.</div>
    </div>
  );
}
