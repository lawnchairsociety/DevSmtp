using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class RcptHandler : ICommandHandler<Rcpt, RcptResult>
    {
        private readonly IDataStore _dataStore;

        public RcptHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<RcptResult> ExecuteAsync(Rcpt command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
