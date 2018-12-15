  import { Reducer } from 'redux';
  import { RootAction } from '../store.js';

  import {
   // STATE_UPDATE,
    GET_ENTITIES,
   // EntityAction
  } from '../actions/entity.js';

  export interface MHState {
    entities: EntitiesState;
    error: string;
  }
  export interface EntitiesState {
    [index: string]: EntityState;
  }

  export interface EntityState {
    id: string;
    state: string;
    attributes: string;
  }

  const INITIAL_STATE: MHState = {
    entities: {},
    error: ''
  };

  const my_home: Reducer<MHState, RootAction> = (state = INITIAL_STATE, action) => {
    switch (action.type) {
      case GET_ENTITIES:
        return {
          ...state,
          entities: action.entities
        };
      default:
        return state;
    }
  };

  // const entities = (state: EntitiesState, action: EntityAction) => {
  //   switch (action.type) {
  //     case STATE_UPDATE:
  //       const entityId = action.entityId;
  //       return {
  //         ...state,
  //         [entityId]: entity(state[entityId], action)
  //       };
  //     default:
  //       return state;
  //   }
  // };
  // const entity = (state: EntityState, action: EntityAction) => {
  //   switch (action.type) {
  //     case STATE_UPDATE:
  //       return {
  //         ...state,
  //         state: state.state
  //       };
  //     default:
  //       return state;
  //   }
  // };

  export default my_home;
  