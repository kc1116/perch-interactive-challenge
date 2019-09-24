import * as React from "react"
import { Header, List, Segment, Visibility} from 'semantic-ui-react'
import Infinite from 'react-infinite';
const style = {
    h1: {
        marginTop: '3em',
    },
    h2: {
        margin: '4em 0em 2em',
    },
    h3: {
        marginTop: '2em',
        padding: '2em 0em',
    },
    last: {
        marginBottom: '300px',
    },
}

class Interaction extends React.Component {
    render() {
        return (
            <List.Item key={Math.random().toString()}>
                <List.Content>
                    <List.Header>
                        <span>{this.props.productName}</span>
                        <span>Published at {this.props.timestamp}</span>
                    </List.Header>
                    <List.Description>
                        <span>Interaction Type: {this.props.interactionType}</span>
                        <span>ID: {this.props.id}</span>
                    </List.Description>
                </List.Content>
            </List.Item>
        );
    }
}

export class EventFeed extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            ws: new WebSocket('ws://localhost:8000/ws'),
            connected: false,
            isInfiniteLoading: false,
            interactions: []
        };
    }

    // handleInfiniteLoad() {
    //     this.setState({
    //         isInfiniteLoading: true
    //     });
    //     setTimeout(function() {
    //         var elemLength = that.state.elements.length,
    //             newElements = that.buildElements(elemLength, elemLength + 1000);
    //         that.setState({
    //             isInfiniteLoading: false,
    //             elements: that.state.elements.concat(newElements)
    //         });
    //     }, 2500);
    // }

    // elementInfiniteLoad() {
    //     return <div className="infinite-list-item">
    //         Loading...
    //     </div>;
    // }

    componentDidMount() {
        this.state.ws.onopen = () => {
            // on connecting, do nothing but log it to the console
            const connected = true;
            this.setState({connected})
        };

        this.state.ws.onmessage = evt => {
            // listen to data sent from the websocket server
            const message = JSON.parse(evt.data);
            let interactions = this.state.interactions;
            interactions.push(<Interaction
                id={message.identifier}
                productName={message.productName}
                interactionType={message.interactionType}
                timestamp={message.timestamp} />);
            this.setState({interactions});
            console.log(message);
        };

        this.state.ws.onclose = () => {
            const connected = false;
            this.setState({connected})
            // automatically try to reconnect on connection loss
        }
    }

    render() {
        return (
            <div>
            <Header as='h2' textAlign='center' style={style.h2} content='Perch Device Interactions' />
            <Header as='h3' textAlign='center' style={style.h3}>
                <span>Websocket status:</span>
                <span style={{'padding': '10px', 'color': this.state.connected ? 'green' : 'red'}}>{this.state.connected ? 'connected' : 'disconnected'}</span>
            </Header>
                <Segment inverted style={{'height': '100vh'}}>
                    <Visibility
                        as={List}
                        continuous={false}
                        once={false}
                        divided inverted relaxed
                    >
                            {this.state.interactions}
                    </Visibility>
                </Segment>
            </div>
        )
    }
}