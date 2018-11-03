import { Reducer } from 'redux';

import {
    STATE_UPDATE,
    EntityAction
  } from '../actions/entity.js';

export interface MHState {
    entities: EntitiesState;
    error: string;
}
export interface EntitiesState {
    [index:string]: EntityState;
  }

  export interface EntityState {
    id: number;
    title: string;
    price: number;
    inventory: number;
  }

  const entity = (state: EntityState, action: EntityAction) => {
    switch (action.type) {
      case STATE_UPDATE:
        return {
          ...state,
          inventory: state.inventory - 1
        };
      default:
        return state;
    }
  };