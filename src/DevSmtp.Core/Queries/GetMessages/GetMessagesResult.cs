using DevSmtp.Core.Models;

namespace DevSmtp.Core.Queries
{
    public sealed class GetMessagesResult : QueryResult
    {
        public GetMessagesResult(IEnumerable<Message> results)
        {
            this.Messages = results;
        }

        public GetMessagesResult(Exception error)
            : base(error)
        {
            this.Messages = new List<Message>();
        }

        public IEnumerable<Message> Messages { get; }
    }
}
