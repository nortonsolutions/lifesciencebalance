// Pure JS ES6 React Timer (no JSX, no Babel)
// export class Timer extends React.Component {

export function Timer (props) {

    var [state,setState] = React.useState({
        timePassed: 0
    });

    React.useEffect(() => {
        props.eventDepot.addListener('tick', (t) => {
            setState({
                timePassed: t
            })
        })
    })

    var minutes = (totalTime) => {
        return Math.floor(totalTime/60);
    }

    var seconds = (totalTime) => {
        return (totalTime%60).toLocaleString("en-US", { minimumIntegerDigits: 2 })
    }

        return (
            React.createElement("div", {
                id: "timerContainer",
                className: "mx-0"
              }, String(minutes(state.timePassed) + ":" + seconds(state.timePassed)), "/", 
                props.timeLimit == Infinity ?
                React.createElement("span", { className: "fas fa-infinity"}) : 
                String(minutes(props.timeLimit) + ":" + seconds(props.timeLimit))
            )
        )

}